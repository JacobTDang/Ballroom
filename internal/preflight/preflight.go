// Package preflight runs the "is everything set up" checks shown on the
// boot screen: Docker running, the practice image built, Ollama
// reachable, the tutor model pulled. Checks are informational, not
// blocking — a failing check surfaces a fix hint but doesn't stop you
// from continuing (something like Docker starting mid-session is common
// enough that hard-blocking would just be annoying).
package preflight

import (
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

// CheckImage reports whether the given Docker image has been built.
func CheckImage(image string) Check {
	cmd := fmt.Sprintf(`docker image inspect %s --format "{{.Id}}"`, image)
	out, err := exec.Command("docker", "image", "inspect", image, "--format", "{{.Id}}").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return Check{
			Name:    CheckNameImage,
			OK:      false,
			Detail:  fmt.Sprintf("%q not built — docker build -f docker/Dockerfile -t %s .", image, image),
			Command: cmd,
			Output:  string(out),
		}
	}
	return Check{Name: CheckNameImage, OK: true, Detail: "built", Command: cmd, Output: string(out)}
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
