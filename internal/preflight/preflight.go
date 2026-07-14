// Package preflight runs the "is everything set up" checks shown on the
// boot screen: Docker running, the practice image built, Ollama
// reachable, the tutor model pulled. Checks are informational, not
// blocking — a failing check surfaces a fix hint but doesn't stop you
// from continuing (something like Docker starting mid-session is common
// enough that hard-blocking would just be annoying).
package preflight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const (
	CheckNameDocker = "Docker daemon"
	CheckNameImage  = "Practice image"
	CheckNameOllama = "Ollama"
	CheckNameModel  = "Tutor model"
)

// Check is the result of one preflight check. Command is the actual
// command/request the check ran, and Output is its real, raw output
// (command stdout/stderr, or HTTP response body) — shown on the boot
// screen so it's transparent what's actually happening, not just a
// checkmark and a canned summary string.
type Check struct {
	Name    string
	OK      bool
	Detail  string
	Command string
	Output  string
}

// CheckDocker reports whether the Docker daemon is reachable.
func CheckDocker() Check {
	const cmd = "docker info"
	out, err := exec.Command("docker", "info").CombinedOutput()
	if err != nil {
		return Check{Name: CheckNameDocker, OK: false, Detail: "not running — start Docker Desktop", Command: cmd, Output: string(out)}
	}
	return Check{Name: CheckNameDocker, OK: true, Detail: "running", Command: cmd, Output: string(out)}
}

// CheckOllama reports whether the Ollama endpoint at host is reachable.
func CheckOllama(host string) Check {
	cmd := "GET " + host + "/api/tags"
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil {
		return Check{Name: CheckNameOllama, OK: false, Detail: "unreachable at " + host + " — is `ollama serve` running?", Command: cmd, Output: err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Check{Name: CheckNameOllama, OK: false, Detail: "unreachable at " + host + " — is `ollama serve` running?", Command: cmd, Output: string(body)}
	}
	return Check{Name: CheckNameOllama, OK: true, Detail: "reachable", Command: cmd, Output: string(body)}
}

// CheckModel reports whether model has been pulled into the Ollama at host.
func CheckModel(host, model string) Check {
	cmd := "GET " + host + "/api/tags"
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil {
		return Check{Name: CheckNameModel, OK: false, Detail: "can't reach Ollama to check", Command: cmd, Output: err.Error()}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Check{Name: CheckNameModel, OK: false, Detail: "can't read Ollama response", Command: cmd, Output: err.Error()}
	}
	if !modelPresent(body, model) {
		return Check{Name: CheckNameModel, OK: false, Detail: model + " not pulled — ollama pull " + model, Command: cmd, Output: string(body)}
	}
	return Check{Name: CheckNameModel, OK: true, Detail: model + " ready", Command: cmd, Output: string(body)}
}

// ListModels queries the Ollama /api/tags endpoint at host and returns the
// names of every locally pulled model, in the order Ollama reports them.
// Unlike CheckModel (an informational boot check that treats failure as
// just another Check result), ListModels is used to populate a picker the
// user actively interacts with, so it fails loud on any network, HTTP, or
// parse error instead of silently returning an empty list.
func ListModels(host string) ([]string, error) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("preflight: list models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("preflight: list models: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("preflight: list models: unexpected status %d", resp.StatusCode)
	}

	names, err := parseModelNames(body)
	if err != nil {
		return nil, fmt.Errorf("preflight: list models: parse response: %w", err)
	}
	return names, nil
}

// parseModelNames extracts the "name" field of every entry in an Ollama
// /api/tags response body.
func parseModelNames(tagsJSON []byte) ([]string, error) {
	var parsed struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(tagsJSON, &parsed); err != nil {
		return nil, err
	}
	names := make([]string, len(parsed.Models))
	for i, m := range parsed.Models {
		names[i] = m.Name
	}
	return names, nil
}

// modelPresent reports whether model appears in an Ollama /api/tags
// response body. Malformed JSON is treated as "not present", not an error
// — this only ever feeds an informational check.
func modelPresent(tagsJSON []byte, model string) bool {
	names, err := parseModelNames(tagsJSON)
	if err != nil {
		return false
	}
	for _, name := range names {
		if name == model {
			return true
		}
	}
	return false
}

// pullStatusLine is one line of Ollama's streamed newline-delimited JSON
// response from POST /api/pull.
type pullStatusLine struct {
	Status    string `json:"status"`
	Error     string `json:"error"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
}

// formatPullStatus turns one pullStatusLine into a single human-readable
// line for display, appending a rough completion percentage when Ollama
// has reported a total size for the current step (it doesn't for every
// status, e.g. "pulling manifest" has no size yet).
func formatPullStatus(s pullStatusLine) string {
	if s.Total > 0 {
		pct := float64(s.Completed) / float64(s.Total) * 100
		return fmt.Sprintf("%s (%.0f%%)", s.Status, pct)
	}
	return s.Status
}

// PullModel streams `ollama pull <model>` via Ollama's HTTP POST
// /api/pull endpoint — not the ollama CLI binary, matching every other
// Ollama interaction in this codebase (CheckModel, ListModels): this runs
// on the host, not inside the practice container, but there's still no
// guarantee the ollama CLI itself is on PATH there, whereas the HTTP API
// is exactly what preflight already depends on being reachable. Streams
// one formatted progress line per response line on the returned channel
// (closed once the response ends), then sends exactly one final result
// (nil on success) on the error channel — same shape as
// orchestrator.BuildImage, so callers can reuse the same "drain lines,
// then read one error" pattern for both.
func PullModel(host, model string) (<-chan string, <-chan error) {
	lineCh := make(chan string, 200)
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		defer close(lineCh)

		payload, err := json.Marshal(map[string]string{"model": model})
		if err != nil {
			errCh <- fmt.Errorf("preflight: pull model: build request: %w", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, host+"/api/pull", bytes.NewReader(payload))
		if err != nil {
			errCh <- fmt.Errorf("preflight: pull model: build request: %w", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// No client timeout: a real model pull can take many minutes on a
		// slow connection, unlike every other Ollama call in this package.
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("preflight: pull model: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errCh <- fmt.Errorf("preflight: pull model: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			raw := strings.TrimSpace(scanner.Text())
			if raw == "" {
				continue
			}
			var status pullStatusLine
			if err := json.Unmarshal([]byte(raw), &status); err != nil {
				continue // skip a line we can't parse rather than failing the whole pull
			}
			if status.Error != "" {
				errCh <- fmt.Errorf("preflight: pull model: %s", status.Error)
				return
			}
			lineCh <- formatPullStatus(status)
		}
		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("preflight: pull model: read response: %w", err)
			return
		}

		errCh <- nil
	}()

	return lineCh, errCh
}
