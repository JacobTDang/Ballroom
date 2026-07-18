package session

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	if err := Reference(cfg, &out); err != nil {
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

	err := Reference(cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected a timeout error when the reference is never revealed, got nil")
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
