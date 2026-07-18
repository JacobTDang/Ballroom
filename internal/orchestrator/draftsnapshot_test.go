package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/draft"
)

const testSnapshotInterval = 5 * time.Millisecond

func TestSnapshotLoop_CapturesLatestWriteBeforeCancel(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	exerciseID := "loop-test"

	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("v1"), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- SnapshotLoop(ctx, dataDir, exerciseID, workspace, testSnapshotInterval)
	}()

	// Let at least one tick fire and capture v1.
	time.Sleep(30 * time.Millisecond)

	// Mid-flight edit, simulating the user saving again in the editor
	// while the session is still running.
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("v2 final"), 0o644); err != nil {
		t.Fatalf("edit workspace: %v", err)
	}

	// Give the loop another tick to catch the edit before cancelling.
	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("SnapshotLoop: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("SnapshotLoop did not return after cancel")
	}

	got, err := os.ReadFile(filepath.Join(draft.Dir(dataDir, exerciseID), "solution.go"))
	if err != nil {
		t.Fatalf("expected draft file to exist: %v", err)
	}
	if string(got) != "v2 final" {
		t.Errorf("draft content = %q, want %q", got, "v2 final")
	}
}

func TestSnapshotLoop_ReturnsPromptlyWhenAlreadyCancelled(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	start := time.Now()
	err := SnapshotLoop(ctx, dataDir, "loop-test", workspace, testSnapshotInterval)
	if err != nil {
		t.Fatalf("SnapshotLoop: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("SnapshotLoop took %v to return after cancellation, want prompt return", elapsed)
	}
}

func TestSnapshotLoop_NeverTicksIsHarmlessNoOp(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir() // no solution.* file at all

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- SnapshotLoop(ctx, dataDir, "empty-workspace", workspace, testSnapshotInterval)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("SnapshotLoop: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("SnapshotLoop did not return after cancel")
	}

	if _, err := os.Stat(draft.Dir(dataDir, "empty-workspace")); !os.IsNotExist(err) {
		t.Errorf("expected no draft dir created when the workspace never had a solution file, stat err = %v", err)
	}
}
