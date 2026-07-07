package session

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

const testPollInterval = 5 * time.Millisecond

// simulateHostReveal mimics what orchestrator.WaitAndReveal does on the host
// side: waits for submit.request, then writes tests.ready. Kept independent
// of the orchestrator package so session's tests don't couple to it.
func simulateHostReveal(t *testing.T, controlDir string) {
	t.Helper()
	requestPath := filepath.Join(controlDir, "submit.request")
	go func() {
		for {
			if _, err := os.Stat(requestPath); err == nil {
				os.WriteFile(filepath.Join(controlDir, "tests.ready"), nil, 0o644)
				return
			}
			time.Sleep(testPollInterval)
		}
	}()
}

func baseConfig(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()
	return Config{
		ControlDir:    filepath.Join(dir, "control"),
		WorkspaceDir:  filepath.Join(dir, "workspace"),
		ExerciseID:    "two-pointers-01",
		Category:      "pattern",
		Language:      "go",
		StartedAt:     time.Now(),
		DBPath:        filepath.Join(dir, "tracker.db"),
		PollInterval:  testPollInterval,
		RevealTimeout: time.Second,
	}
}

func mkdirs(t *testing.T, cfg Config) {
	t.Helper()
	if err := os.MkdirAll(cfg.ControlDir, 0o755); err != nil {
		t.Fatalf("mkdir control: %v", err)
	}
	if err := os.MkdirAll(cfg.WorkspaceDir, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
}

func TestSubmit_PassResult(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want %q", attempt.Result, tracker.ResultPass)
	}
}

func TestSubmit_FailResult(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "false"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultFail {
		t.Errorf("Result = %q, want %q", attempt.Result, tracker.ResultFail)
	}
}

func TestSubmit_ComputesTimeSpent(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	cfg.StartedAt = time.Now().Add(-10 * time.Minute)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.TimeSpentMin < 9.9 || attempt.TimeSpentMin > 10.5 {
		t.Errorf("TimeSpentMin = %v, want ~10", attempt.TimeSpentMin)
	}
}

func TestSubmit_PromptsForAndStoresNotes(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("great problem\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Notes != "great problem" {
		t.Errorf("Notes = %q, want %q", attempt.Notes, "great problem")
	}
}

func TestSubmit_TimesOutIfTestsNeverRevealed(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	cfg.RevealTimeout = 30 * time.Millisecond
	mkdirs(t, cfg)
	// deliberately no simulateHostReveal call

	_, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected timeout error when tests are never revealed, got nil")
	}
}

func TestSubmit_WritesAttemptToTrackerDB(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("tracker.Open: %v", err)
	}
	defer tr.Close()

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt logged, got %d", len(attempts))
	}
	if attempts[0].ID != attempt.ID {
		t.Errorf("logged attempt ID = %d, want %d", attempts[0].ID, attempt.ID)
	}
	if attempts[0].ExerciseID != cfg.ExerciseID {
		t.Errorf("logged ExerciseID = %q, want %q", attempts[0].ExerciseID, cfg.ExerciseID)
	}
}
