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

// Check is the result of one preflight check.
type Check struct {
	Name   string
	OK     bool
	Detail string
}

// CheckDocker reports whether the Docker daemon is reachable.
func CheckDocker() Check {
	if err := exec.Command("docker", "info").Run(); err != nil {
		return Check{Name: CheckNameDocker, OK: false, Detail: "not running — start Docker Desktop"}
	}
	return Check{Name: CheckNameDocker, OK: true, Detail: "running"}
}

// CheckImage reports whether the given Docker image has been built.
func CheckImage(image string) Check {
	out, err := exec.Command("docker", "image", "inspect", image, "--format", "{{.Id}}").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return Check{
			Name:   CheckNameImage,
			OK:     false,
			Detail: fmt.Sprintf("%q not built — docker build -f docker/Dockerfile -t %s .", image, image),
		}
	}
	return Check{Name: CheckNameImage, OK: true, Detail: "built"}
}

// CheckOllama reports whether the Ollama endpoint at host is reachable.
func CheckOllama(host string) Check {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil || resp.StatusCode != http.StatusOK {
		return Check{Name: CheckNameOllama, OK: false, Detail: "unreachable at " + host + " — is `ollama serve` running?"}
	}
	defer resp.Body.Close()
	return Check{Name: CheckNameOllama, OK: true, Detail: "reachable"}
}

// CheckModel reports whether model has been pulled into the Ollama at host.
func CheckModel(host, model string) Check {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil {
		return Check{Name: CheckNameModel, OK: false, Detail: "can't reach Ollama to check"}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Check{Name: CheckNameModel, OK: false, Detail: "can't read Ollama response"}
	}
	if !modelPresent(body, model) {
		return Check{Name: CheckNameModel, OK: false, Detail: model + " not pulled — ollama pull " + model}
	}
	return Check{Name: CheckNameModel, OK: true, Detail: model + " ready"}
}

// modelPresent reports whether model appears in an Ollama /api/tags
// response body. Malformed JSON is treated as "not present", not an error
// — this only ever feeds an informational check.
func modelPresent(tagsJSON []byte, model string) bool {
	var parsed struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(tagsJSON, &parsed); err != nil {
		return false
	}
	for _, m := range parsed.Models {
		if m.Name == model {
			return true
		}
	}
	return false
}
