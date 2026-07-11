// Package session implements the in-container side of `practice submit`:
// request the hidden tests be revealed, run test_command once they land,
// and log the result. Runs inside the practice container, talking to the
// host-side orchestrator only through files under a shared control dir
// (bind-mounted from the host) — no network, no Docker socket.
package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// lastTestResultFile is the well-known dotfile written into a workspace
// after every submit, holding the raw test output the tutor's
// read_test_output tool reads from (internal/tutor). Living directly in
// the workspace mirrors how the tutor already finds solution.*/problem.md
// there — no DB schema change, and it's naturally scoped to just this
// exercise run.
const lastTestResultFile = ".ballroom-last-test-result.json"

// lastTestResult is the JSON shape written to lastTestResultFile.
type lastTestResult struct {
	Result      string    `json:"result"`
	Output      string    `json:"output"`
	TestCommand string    `json:"test_command"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// Config describes one submit invocation. All paths are as seen from
// inside the container.
type Config struct {
	ControlDir    string
	WorkspaceDir  string
	TestCommand   string
	ExerciseID    string
	Category      string
	Language      string
	StartedAt     time.Time
	DBPath        string
	PollInterval  time.Duration
	RevealTimeout time.Duration
}

// Submit requests the hidden tests, waits for them to be revealed, runs
// test_command, logs the attempt, and returns it.
func Submit(cfg Config, stdin io.Reader, stdout io.Writer) (tracker.Attempt, error) {
	if err := requestReveal(cfg); err != nil {
		return tracker.Attempt{}, err
	}
	if err := waitForReady(cfg); err != nil {
		return tracker.Attempt{}, err
	}

	result, output := runTestCommand(cfg)
	fmt.Fprintf(stdout, "\nresult: %s\n%s\n", result, output)

	if err := writeLastTestResult(cfg, result, output); err != nil {
		fmt.Fprintf(stdout, "warning: could not save test output for the tutor: %v\n", err)
		// Not fatal — same graceful-degradation philosophy as the rest of
		// the tutor-adjacent code. A submission should still get logged
		// to the tracker DB even if this write fails.
	}

	notes := promptNotes(stdin, stdout)

	attempt := tracker.Attempt{
		ExerciseID:   cfg.ExerciseID,
		Category:     cfg.Category,
		Language:     cfg.Language,
		Date:         time.Now().Format("2006-01-02"),
		TimeSpentMin: time.Since(cfg.StartedAt).Minutes(),
		Result:       result,
		Notes:        notes,
	}

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: open tracker: %w", err)
	}
	defer tr.Close()

	id, err := tr.LogAttempt(attempt)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: log attempt: %w", err)
	}
	attempt.ID = id
	return attempt, nil
}

func requestReveal(cfg Config) error {
	path := filepath.Join(cfg.ControlDir, "submit.request")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		return fmt.Errorf("session: write submit request: %w", err)
	}
	return nil
}

func waitForReady(cfg Config) error {
	readyPath := filepath.Join(cfg.ControlDir, "tests.ready")
	deadline := time.Now().Add(cfg.RevealTimeout)

	for time.Now().Before(deadline) {
		if _, err := os.Stat(readyPath); err == nil {
			return nil
		}
		time.Sleep(cfg.PollInterval)
	}
	return fmt.Errorf("session: timed out after %s waiting for hidden tests to be revealed", cfg.RevealTimeout)
}

func runTestCommand(cfg Config) (result string, output string) {
	cmd := exec.Command("sh", "-c", cfg.TestCommand)
	cmd.Dir = cfg.WorkspaceDir
	out, err := cmd.CombinedOutput()

	result = tracker.ResultPass
	if err != nil {
		result = tracker.ResultFail
	}
	return result, string(out)
}

// writeLastTestResult persists result/output to lastTestResultFile in
// the workspace so the tutor's read_test_output tool has something to
// read after this submit — see internal/tutor for the reader.
func writeLastTestResult(cfg Config, result, output string) error {
	data, err := json.Marshal(lastTestResult{
		Result:      result,
		Output:      output,
		TestCommand: cfg.TestCommand,
		RecordedAt:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("session: marshal last test result: %w", err)
	}

	path := filepath.Join(cfg.WorkspaceDir, lastTestResultFile)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("session: write last test result: %w", err)
	}
	return nil
}

func promptNotes(stdin io.Reader, stdout io.Writer) string {
	fmt.Fprint(stdout, "Notes (optional): ")
	scanner := bufio.NewScanner(stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
