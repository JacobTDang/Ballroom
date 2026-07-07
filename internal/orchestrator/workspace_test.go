package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareWorkspace_CopiesRepoContents(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo)
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()

	if workspace == repo {
		t.Fatal("workspace must be a distinct directory from the source repo")
	}

	got, err := os.ReadFile(filepath.Join(workspace, "solution.go"))
	if err != nil {
		t.Fatalf("expected solution.go copied into workspace: %v", err)
	}
	if string(got) != "package main" {
		t.Errorf("copied content = %q, want %q", got, "package main")
	}
}

func TestPrepareWorkspace_CleanupRemovesDir(t *testing.T) {
	repo := t.TempDir()
	os.WriteFile(filepath.Join(repo, "f.go"), []byte("x"), 0o644)

	workspace, cleanup, err := PrepareWorkspace(repo)
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}

	cleanup()

	if _, err := os.Stat(workspace); !os.IsNotExist(err) {
		t.Errorf("expected workspace dir removed after cleanup, stat err = %v", err)
	}
}

func TestPrepareWorkspace_SourceRepoUnaffectedByWorkspaceEdits(t *testing.T) {
	repo := t.TempDir()
	original := "package main\n\nfunc TwoSum() {}\n"
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte(original), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo)
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()

	// Simulate what happens during a session: edits in the workspace,
	// and a hidden test getting revealed into it.
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("mutated"), 0o644); err != nil {
		t.Fatalf("mutate workspace copy: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution_test.go"), []byte("hidden"), 0o644); err != nil {
		t.Fatalf("write revealed test into workspace: %v", err)
	}

	gotSource, err := os.ReadFile(filepath.Join(repo, "solution.go"))
	if err != nil {
		t.Fatalf("read source repo file: %v", err)
	}
	if string(gotSource) != original {
		t.Errorf("source repo solution.go was mutated: got %q, want unchanged %q", gotSource, original)
	}
	if _, err := os.Stat(filepath.Join(repo, "solution_test.go")); !os.IsNotExist(err) {
		t.Error("hidden test leaked into source repo — this is exactly the bug this fix is for")
	}
}
