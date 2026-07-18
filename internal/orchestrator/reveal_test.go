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

// seedReferenceFixture creates referenceSrc/solution.py (the "reference
// solution") and testsSrc/hidden_test.py (the "hidden tests" a reference
// must satisfy), mirroring one real exercise's .reference/ + tests/<id>/
// layout.
func seedReferenceFixture(t *testing.T, referenceSrc, testsSrc string) {
	t.Helper()
	if err := os.MkdirAll(referenceSrc, 0o755); err != nil {
		t.Fatalf("mkdir referenceSrc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(referenceSrc, "solution.py"), []byte("def solve(): return 42"), 0o644); err != nil {
		t.Fatalf("seed reference solution: %v", err)
	}
	if err := os.MkdirAll(testsSrc, 0o755); err != nil {
		t.Fatalf("mkdir testsSrc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testsSrc, "hidden_test.py"), []byte("def test_solve(): assert solve() == 42"), 0o644); err != nil {
		t.Fatalf("seed hidden test: %v", err)
	}
}

func TestWaitAndRevealReference_HidesUntilRequested(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()
	seedReferenceFixture(t, referenceSrc, testsSrc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval)
	}()

	time.Sleep(20 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(workspace, "reference")); err == nil {
		t.Fatal("reference/ appeared in workspace before it was requested")
	}

	if err := os.WriteFile(filepath.Join(controlDir, "reference.request"), nil, 0o644); err != nil {
		t.Fatalf("write reference.request: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("WaitAndRevealReference: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WaitAndRevealReference did not return after reference.request was written")
	}

	if _, err := os.Stat(filepath.Join(controlDir, "reference.ready")); err != nil {
		t.Errorf("expected reference.ready marker to exist: %v", err)
	}
}

// TestWaitAndRevealReference_LandsInSubdirectoryUserSolutionUntouched is
// the actual hiding/non-clobbering guarantee: the reference solution
// must never land at the workspace root next to (or instead of) the
// user's own in-progress solution file, only inside reference/.
func TestWaitAndRevealReference_LandsInSubdirectoryUserSolutionUntouched(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()
	seedReferenceFixture(t, referenceSrc, testsSrc)

	const usersOwnCode = "def solve(): return 'still working on it'"
	if err := os.WriteFile(filepath.Join(workspace, "solution.py"), []byte(usersOwnCode), 0o644); err != nil {
		t.Fatalf("seed user's own solution: %v", err)
	}

	if err := os.WriteFile(filepath.Join(controlDir, "reference.request"), nil, 0o644); err != nil {
		t.Fatalf("write reference.request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval); err != nil {
		t.Fatalf("WaitAndRevealReference: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(workspace, "solution.py"))
	if err != nil {
		t.Fatalf("read user's solution.py: %v", err)
	}
	if string(got) != usersOwnCode {
		t.Errorf("user's solution.py = %q, want untouched %q", got, usersOwnCode)
	}

	refSolution, err := os.ReadFile(filepath.Join(workspace, "reference", "solution.py"))
	if err != nil {
		t.Fatalf("read reference/solution.py: %v", err)
	}
	if string(refSolution) != "def solve(): return 42" {
		t.Errorf("reference/solution.py = %q, want the reference content", refSolution)
	}

	// The hidden tests the reference must satisfy are revealed alongside
	// it -- reading the answer without the tests it satisfies is half a
	// lesson.
	if _, err := os.Stat(filepath.Join(workspace, "reference", "hidden_test.py")); err != nil {
		t.Errorf("expected hidden tests copied into reference/ too: %v", err)
	}
}

func TestWaitAndRevealReference_MissingReferenceDirGivesClearErrorPromptly(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), "does-not-exist", ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()

	if err := os.WriteFile(filepath.Join(controlDir, "reference.request"), nil, 0o644); err != nil {
		t.Fatalf("write reference.request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	err := WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error for a missing .reference dir, got nil")
	}
	if elapsed > time.Second {
		t.Errorf("WaitAndRevealReference took %v to error out on a missing reference dir, want a prompt failure rather than hanging until ctx times out", elapsed)
	}
	if _, statErr := os.Stat(filepath.Join(controlDir, "reference.ready")); statErr == nil {
		t.Error("reference.ready was written despite the reveal failing")
	}
}

func TestWaitAndRevealReference_ContextCancelReturnsPromptly(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := t.TempDir()
	testsSrc := t.TempDir()
	workspace := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	start := time.Now()
	err := WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval)
	if err == nil {
		t.Fatal("expected error when context is already cancelled, got nil")
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("WaitAndRevealReference took %v to return after cancellation, want prompt return", elapsed)
	}
}

// TestWaitAndRevealReference_RepeatedRequestWritesAreIdempotent covers a
// request file rewritten (e.g. a second `ballroom reference` press)
// before the watcher's next poll tick: it must still complete cleanly
// with the reference correctly revealed exactly once, not error or
// double up.
func TestWaitAndRevealReference_RepeatedRequestWritesAreIdempotent(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()
	seedReferenceFixture(t, referenceSrc, testsSrc)

	requestPath := filepath.Join(controlDir, "reference.request")
	if err := os.WriteFile(requestPath, nil, 0o644); err != nil {
		t.Fatalf("write reference.request (1st): %v", err)
	}
	// Rewrite it again right away, before the watcher below has even
	// started polling.
	if err := os.WriteFile(requestPath, nil, 0o644); err != nil {
		t.Fatalf("write reference.request (2nd): %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval); err != nil {
		t.Fatalf("WaitAndRevealReference: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(workspace, "reference", "solution.py"))
	if err != nil {
		t.Fatalf("read reference/solution.py: %v", err)
	}
	if string(got) != "def solve(): return 42" {
		t.Errorf("reference/solution.py = %q, want the reference content", got)
	}
}

// TestReveal_ReferenceThenSubmit and TestReveal_SubmitThenReference cover
// the generalized watcher handling both request types in either order —
// each one watched and fulfilled independently, so a session can freely
// ask for the reference before, after, or without ever submitting.
func TestReveal_ReferenceThenSubmit(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()
	seedReferenceFixture(t, referenceSrc, testsSrc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	refDone := make(chan error, 1)
	go func() {
		refDone <- WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval)
	}()
	submitDone := make(chan error, 1)
	go func() {
		submitDone <- WaitAndReveal(ctx, controlDir, testsSrc, workspace, testPollInterval)
	}()

	if err := os.WriteFile(filepath.Join(controlDir, "reference.request"), nil, 0o644); err != nil {
		t.Fatalf("write reference.request: %v", err)
	}
	if err := <-refDone; err != nil {
		t.Fatalf("WaitAndRevealReference: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "reference", "solution.py")); err != nil {
		t.Errorf("expected reference revealed first: %v", err)
	}

	// Submit hasn't been requested yet -- hidden tests must still be
	// absent from the workspace root.
	if _, err := os.Stat(filepath.Join(workspace, "hidden_test.py")); err == nil {
		t.Fatal("hidden tests appeared at workspace root before submit was requested")
	}

	if err := os.WriteFile(filepath.Join(controlDir, "submit.request"), nil, 0o644); err != nil {
		t.Fatalf("write submit.request: %v", err)
	}
	if err := <-submitDone; err != nil {
		t.Fatalf("WaitAndReveal: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "hidden_test.py")); err != nil {
		t.Errorf("expected hidden tests revealed at workspace root after submit: %v", err)
	}
}

func TestReveal_SubmitThenReference(t *testing.T) {
	controlDir := t.TempDir()
	referenceSrc := filepath.Join(t.TempDir(), ".reference")
	testsSrc := t.TempDir()
	workspace := t.TempDir()
	seedReferenceFixture(t, referenceSrc, testsSrc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	refDone := make(chan error, 1)
	go func() {
		refDone <- WaitAndRevealReference(ctx, controlDir, referenceSrc, testsSrc, workspace, testPollInterval)
	}()
	submitDone := make(chan error, 1)
	go func() {
		submitDone <- WaitAndReveal(ctx, controlDir, testsSrc, workspace, testPollInterval)
	}()

	if err := os.WriteFile(filepath.Join(controlDir, "submit.request"), nil, 0o644); err != nil {
		t.Fatalf("write submit.request: %v", err)
	}
	if err := <-submitDone; err != nil {
		t.Fatalf("WaitAndReveal: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "hidden_test.py")); err != nil {
		t.Errorf("expected hidden tests revealed at workspace root after submit: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "reference")); err == nil {
		t.Fatal("reference/ appeared before it was ever requested")
	}

	if err := os.WriteFile(filepath.Join(controlDir, "reference.request"), nil, 0o644); err != nil {
		t.Fatalf("write reference.request: %v", err)
	}
	if err := <-refDone; err != nil {
		t.Fatalf("WaitAndRevealReference: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "reference", "solution.py")); err != nil {
		t.Errorf("expected reference revealed after submit: %v", err)
	}
}
