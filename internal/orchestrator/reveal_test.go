package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testPollInterval = 5 * time.Millisecond

func TestWaitAndReveal_HidesUntilRequested(t *testing.T) {
	controlDir := t.TempDir()
	testsSrc := t.TempDir()
	workspace := t.TempDir()

	if err := os.WriteFile(filepath.Join(testsSrc, "hidden_test.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed testsSrc: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- WaitAndReveal(ctx, controlDir, testsSrc, workspace, testPollInterval)
	}()

	// Give the watcher a moment to start; the hidden test file must not
	// exist in the workspace yet — this is the actual "hiding" guarantee.
	time.Sleep(20 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(workspace, "hidden_test.go")); err == nil {
		t.Fatal("hidden test file appeared in workspace before submit was requested")
	}

	if err := os.WriteFile(filepath.Join(controlDir, "submit.request"), nil, 0o644); err != nil {
		t.Fatalf("write submit.request: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("WaitAndReveal: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WaitAndReveal did not return after submit.request was written")
	}

	got, err := os.ReadFile(filepath.Join(workspace, "hidden_test.go"))
	if err != nil {
		t.Fatalf("expected hidden_test.go copied into workspace: %v", err)
	}
	if string(got) != "package main" {
		t.Errorf("copied file content = %q, want %q", got, "package main")
	}

	if _, err := os.Stat(filepath.Join(controlDir, "tests.ready")); err != nil {
		t.Errorf("expected tests.ready marker to exist: %v", err)
	}
}

func TestWaitAndReveal_CopiesNestedDirectories(t *testing.T) {
	controlDir := t.TempDir()
	testsSrc := t.TempDir()
	workspace := t.TempDir()

	nested := filepath.Join(testsSrc, "sub", "dir")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nested, "case1_test.go"), []byte("case1"), 0o644); err != nil {
		t.Fatalf("seed nested file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- WaitAndReveal(ctx, controlDir, testsSrc, workspace, testPollInterval) }()

	if err := os.WriteFile(filepath.Join(controlDir, "submit.request"), nil, 0o644); err != nil {
		t.Fatalf("write submit.request: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("WaitAndReveal: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(workspace, "sub", "dir", "case1_test.go"))
	if err != nil {
		t.Fatalf("expected nested file copied into workspace: %v", err)
	}
	if string(got) != "case1" {
		t.Errorf("copied nested file content = %q, want %q", got, "case1")
	}
}

func TestWaitAndReveal_ContextCancelReturnsPromptly(t *testing.T) {
	controlDir := t.TempDir()
	testsSrc := t.TempDir()
	workspace := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	start := time.Now()
	err := WaitAndReveal(ctx, controlDir, testsSrc, workspace, testPollInterval)
	if err == nil {
		t.Fatal("expected error when context is already cancelled, got nil")
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("WaitAndReveal took %v to return after cancellation, want prompt return", elapsed)
	}
}
