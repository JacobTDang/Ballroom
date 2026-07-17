package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareWorkspace_CopiesRepoContents(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo, "")
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

func TestPrepareWorkspace_RendersProblemTextAlongsideProblemMd(t *testing.T) {
	repo := t.TempDir()
	md := "# Two Sum\n\nreturn indices of the **two numbers** that add to `target`\n"
	if err := os.WriteFile(filepath.Join(repo, "problem.md"), []byte(md), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo, "")
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()

	got, err := os.ReadFile(filepath.Join(workspace, "problem.txt"))
	if err != nil {
		t.Fatalf("expected problem.txt rendered into workspace: %v", err)
	}
	text := string(got)
	if strings.Contains(text, "**") || strings.Contains(text, "`") || strings.Contains(text, "#") {
		t.Errorf("problem.txt still contains markdown markers:\n%s", text)
	}
	if !strings.Contains(text, "Two Sum") || !strings.Contains(text, "two numbers") {
		t.Errorf("problem.txt lost real content:\n%s", text)
	}
	// The markdown source must stay in the workspace untouched -- the
	// tutor's read_problem_statement tool reads problem.md.
	if _, err := os.Stat(filepath.Join(workspace, "problem.md")); err != nil {
		t.Errorf("problem.md missing from workspace: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "problem.txt")); !os.IsNotExist(err) {
		t.Error("problem.txt leaked into the source repo -- render must go to the workspace only")
	}
}

func TestPrepareWorkspace_NoProblemTxtWhenExerciseHasNoProblemMd(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo, "")
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()

	if _, err := os.Stat(filepath.Join(workspace, "problem.txt")); !os.IsNotExist(err) {
		t.Errorf("problem.txt exists for an exercise with no problem.md, stat err = %v", err)
	}
}

func TestPrepareWorkspace_CleanupRemovesDir(t *testing.T) {
	repo := t.TempDir()
	os.WriteFile(filepath.Join(repo, "f.go"), []byte("x"), 0o644)

	workspace, cleanup, err := PrepareWorkspace(repo, "")
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

	workspace, cleanup, err := PrepareWorkspace(repo, "")
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

func TestPrepareWorkspace_VideoFooterOnProblemTxt(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "problem.md"), []byte("# T\n\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dir, cleanup, err := PrepareWorkspace(repo, "https://youtu.be/abc123")
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()
	b, err := os.ReadFile(filepath.Join(dir, "problem.txt"))
	if err != nil {
		t.Fatalf("read problem.txt: %v", err)
	}
	if !strings.Contains(string(b), "solution video (spoilers!): https://youtu.be/abc123") {
		t.Errorf("problem.txt missing the video footer:\n%s", b)
	}

	// No URL: no footer.
	dir2, cleanup2, err := PrepareWorkspace(repo, "")
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup2()
	b2, _ := os.ReadFile(filepath.Join(dir2, "problem.txt"))
	if strings.Contains(string(b2), "solution video") {
		t.Errorf("problem.txt has a footer with no URL:\n%s", b2)
	}
}
