package orchestrator

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestPrepareWorkspace_CopiesRepoContents(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo, "", "")
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

	workspace, cleanup, err := PrepareWorkspace(repo, "", "")
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

	workspace, cleanup, err := PrepareWorkspace(repo, "", "")
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

	workspace, cleanup, err := PrepareWorkspace(repo, "", "")
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

	workspace, cleanup, err := PrepareWorkspace(repo, "", "")
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

	dir, cleanup, err := PrepareWorkspace(repo, "https://youtu.be/abc123", "")
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
	dir2, cleanup2, err := PrepareWorkspace(repo, "", "")
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup2()
	b2, _ := os.ReadFile(filepath.Join(dir2, "problem.txt"))
	if strings.Contains(string(b2), "solution video") {
		t.Errorf("problem.txt has a footer with no URL:\n%s", b2)
	}
}

func TestPrepareWorkspace_DraftOverlaysStarterAndStillRendersProblemText(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte("package main\n\n// starter\n"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "problem.md"), []byte("# T\n\nbody\n"), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}

	draftDir := t.TempDir()
	edited := "package main\n\n// the user's in-progress edit\n"
	if err := os.WriteFile(filepath.Join(draftDir, "solution.go"), []byte(edited), 0o644); err != nil {
		t.Fatalf("seed draft: %v", err)
	}

	workspace, cleanup, err := PrepareWorkspace(repo, "", draftDir)
	if err != nil {
		t.Fatalf("PrepareWorkspace: %v", err)
	}
	defer cleanup()

	got, err := os.ReadFile(filepath.Join(workspace, "solution.go"))
	if err != nil {
		t.Fatalf("read workspace solution.go: %v", err)
	}
	if string(got) != edited {
		t.Errorf("solution.go = %q, want the overlaid draft content %q", got, edited)
	}

	// The overlay must not disturb the rendered problem statement.
	if _, err := os.Stat(filepath.Join(workspace, "problem.txt")); err != nil {
		t.Errorf("expected problem.txt still rendered with a draft overlay present: %v", err)
	}
	if _, err := os.Stat(filepath.Join(draftDir, "problem.txt")); !os.IsNotExist(err) {
		t.Error("overlay must not write anything back into the draft dir")
	}
}

// TestPrepareWorkspace_EmptyDraftDirIsByteIdenticalToNoOverlay is a
// regression pin: an empty draftDir ("" -- an exercise never
// practiced before, the overwhelming common case) must produce exactly
// what PrepareWorkspace produced before draft overlay existed at all.
// A draftDir that resolves to a real but never-snapshotted directory
// (draft.Dir's return value on an exercise with no draft yet) must
// behave identically too, since that's what every real call site now
// passes on a fresh exercise.
func TestPrepareWorkspace_EmptyDraftDirIsByteIdenticalToNoOverlay(t *testing.T) {
	repo := t.TempDir()
	starter := "package main\n\nfunc TwoSum() {}\n"
	if err := os.WriteFile(filepath.Join(repo, "solution.go"), []byte(starter), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}
	md := "# Two Sum\n\nbody\n"
	if err := os.WriteFile(filepath.Join(repo, "problem.md"), []byte(md), 0o644); err != nil {
		t.Fatalf("seed repo: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	helper := "package main\n"
	if err := os.WriteFile(filepath.Join(repo, "sub", "helper.go"), []byte(helper), 0o644); err != nil {
		t.Fatal(err)
	}

	withoutDraft, cleanup1, err := PrepareWorkspace(repo, "", "")
	if err != nil {
		t.Fatalf("PrepareWorkspace (empty draftDir): %v", err)
	}
	defer cleanup1()

	neverSnapshotted := filepath.Join(t.TempDir(), "drafts", "some-exercise-id")
	withRealButEmptyDraftDir, cleanup2, err := PrepareWorkspace(repo, "", neverSnapshotted)
	if err != nil {
		t.Fatalf("PrepareWorkspace (real but empty draftDir): %v", err)
	}
	defer cleanup2()

	for _, workspace := range []string{withoutDraft, withRealButEmptyDraftDir} {
		got, err := os.ReadFile(filepath.Join(workspace, "solution.go"))
		if err != nil {
			t.Fatalf("read %s/solution.go: %v", workspace, err)
		}
		if string(got) != starter {
			t.Errorf("%s: solution.go = %q, want unmodified starter %q", workspace, got, starter)
		}

		gotHelper, err := os.ReadFile(filepath.Join(workspace, "sub", "helper.go"))
		if err != nil {
			t.Fatalf("read %s/sub/helper.go: %v", workspace, err)
		}
		if string(gotHelper) != helper {
			t.Errorf("%s: sub/helper.go = %q, want unmodified %q", workspace, gotHelper, helper)
		}

		var names []string
		if err := filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			rel, relErr := filepath.Rel(workspace, path)
			if relErr != nil {
				return relErr
			}
			names = append(names, rel)
			return nil
		}); err != nil {
			t.Fatalf("walk %s: %v", workspace, err)
		}
		sort.Strings(names)
		want := []string{"problem.md", "problem.txt", "solution.go", filepath.Join("sub", "helper.go")}
		sort.Strings(want)
		if !reflect.DeepEqual(names, want) {
			t.Errorf("%s: file listing = %v, want %v (no extra files from the overlay step)", workspace, names, want)
		}
	}
}
