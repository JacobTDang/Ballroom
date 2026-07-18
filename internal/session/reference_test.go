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

// simulateHostReferenceReveal mimics what
// orchestrator.WaitAndRevealReference does on the host side: waits for
// reference.request, then writes reference.ready. Kept independent of
// the orchestrator package for the same reason simulateHostReveal is,
// in submit_test.go.
func simulateHostReferenceReveal(t *testing.T, controlDir string) {
	t.Helper()
	requestPath := filepath.Join(controlDir, "reference.request")
	go func() {
		for {
			if _, err := os.Stat(requestPath); err == nil {
				os.WriteFile(filepath.Join(controlDir, "reference.ready"), nil, 0o644)
				return
			}
			time.Sleep(testPollInterval)
		}
	}()
}

func TestReference_RequestsThenWaitsForReady(t *testing.T) {
	cfg := baseConfig(t)
	mkdirs(t, cfg)
	simulateHostReferenceReveal(t, cfg.ControlDir)

	var out bytes.Buffer
	if _, err := Reference(cfg, &out); err != nil {
		t.Fatalf("Reference: %v", err)
	}

	if _, err := os.Stat(filepath.Join(cfg.ControlDir, "reference.request")); err != nil {
		t.Errorf("expected reference.request written: %v", err)
	}
}

func TestReference_TimesOutIfNeverRevealed(t *testing.T) {
	cfg := baseConfig(t)
	cfg.RevealTimeout = 30 * time.Millisecond
	mkdirs(t, cfg)
	// deliberately no simulateHostReferenceReveal call

	_, err := Reference(cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected a timeout error when the reference is never revealed, got nil")
	}
}

// TestReference_TimeoutNeverRecordsAnAttempt: a reveal that never
// completes must not log anything -- there was no real "give up"
// moment, just a broken handshake.
func TestReference_TimeoutNeverRecordsAnAttempt(t *testing.T) {
	cfg := baseConfig(t)
	cfg.RevealTimeout = 30 * time.Millisecond
	mkdirs(t, cfg)

	if _, err := Reference(cfg, &bytes.Buffer{}); err == nil {
		t.Fatal("expected a timeout error, got nil")
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
	if len(attempts) != 0 {
		t.Errorf("attempts = %+v, want none logged after a reveal timeout", attempts)
	}
}

// TestReference_RecordsGaveUpAttempt is the core of issue #238: asking
// to see the reference IS the outcome, recorded honestly rather than
// left for a submit that may never come.
func TestReference_RecordsGaveUpAttempt(t *testing.T) {
	cfg := baseConfig(t)
	mkdirs(t, cfg)
	simulateHostReferenceReveal(t, cfg.ControlDir)

	attempt, err := Reference(cfg, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Reference: %v", err)
	}
	if attempt.Result != tracker.ResultGaveUp {
		t.Errorf("Result = %q, want %q", attempt.Result, tracker.ResultGaveUp)
	}
	if attempt.ExerciseID != cfg.ExerciseID {
		t.Errorf("ExerciseID = %q, want %q", attempt.ExerciseID, cfg.ExerciseID)
	}
	if attempt.ID == 0 {
		t.Error("expected a non-zero logged attempt ID")
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
	if len(attempts) != 1 || attempts[0].Result != tracker.ResultGaveUp {
		t.Errorf("logged attempts = %+v, want a single gave-up row", attempts)
	}
}

func TestReference_RecordsHintsUsedTutorModeAndModelFromDotfile(t *testing.T) {
	cfg := baseConfig(t)
	mkdirs(t, cfg)
	simulateHostReferenceReveal(t, cfg.ControlDir)
	writeTutorStateFixture(t, cfg.WorkspaceDir, tutorState{HintsUsed: 2, TutorMode: "hints-first", Model: "llama3.1:8b"})

	attempt, err := Reference(cfg, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Reference: %v", err)
	}
	if attempt.HintsUsed == nil || *attempt.HintsUsed != 2 {
		t.Errorf("HintsUsed = %v, want 2", attempt.HintsUsed)
	}
	if attempt.TutorMode == nil || *attempt.TutorMode != "hints-first" {
		t.Errorf("TutorMode = %v, want %q", attempt.TutorMode, "hints-first")
	}
	if attempt.Model == nil || *attempt.Model != "llama3.1:8b" {
		t.Errorf("Model = %v, want %q", attempt.Model, "llama3.1:8b")
	}
}

func TestReference_MissingTutorStateDotfileDegradesToZeroValue(t *testing.T) {
	cfg := baseConfig(t)
	mkdirs(t, cfg)
	simulateHostReferenceReveal(t, cfg.ControlDir)
	// Deliberately no tutor-state dotfile written.

	var out bytes.Buffer
	attempt, err := Reference(cfg, &out)
	if err != nil {
		t.Fatalf("Reference: %v", err)
	}
	if attempt.HintsUsed == nil || *attempt.HintsUsed != 0 {
		t.Errorf("HintsUsed = %v, want a non-nil 0", attempt.HintsUsed)
	}
	if attempt.TutorMode == nil || *attempt.TutorMode != "" {
		t.Errorf("TutorMode = %v, want a non-nil empty string", attempt.TutorMode)
	}
	if strings.Contains(out.String(), "warn") || strings.Contains(out.String(), "error") {
		t.Errorf("output %q should not warn about a missing tutor-state dotfile", out.String())
	}
}

func TestReferenceSolutionPath_FindsRevealedSolutionFile(t *testing.T) {
	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, "reference"), 0o755); err != nil {
		t.Fatalf("mkdir reference: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "reference", "solution.py"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write reference solution: %v", err)
	}

	got, err := ReferenceSolutionPath(workspace)
	if err != nil {
		t.Fatalf("ReferenceSolutionPath: %v", err)
	}
	want := filepath.Join(workspace, "reference", "solution.py")
	if got != want {
		t.Errorf("ReferenceSolutionPath = %q, want %q", got, want)
	}
}

func TestReferenceSolutionPath_MissingReturnsClearError(t *testing.T) {
	workspace := t.TempDir()
	if _, err := ReferenceSolutionPath(workspace); err == nil {
		t.Fatal("expected an error when no reference solution has been revealed yet, got nil")
	}
}
