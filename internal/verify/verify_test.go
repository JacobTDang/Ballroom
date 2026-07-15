package verify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile creates path (and its parent dirs) with contents.
func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

const goModContents = "module practice/fixture\n\ngo 1.22\n"

const goTestContents = `package main

import "testing"

func TestDouble(t *testing.T) {
	if got := Double(3); got != 6 {
		t.Errorf("Double(3) = %d, want 6", got)
	}
}
`

// fixtureExercise builds a minimal but real go-based exercise directory:
// exercise.json, repo/{go.mod, solution.go: starterBody}, and (if
// referenceBody != "") .reference/solution.go: referenceBody. Returns the
// exercise dir and a sibling tests dir.
func fixtureExercise(t *testing.T, starterBody, referenceBody string) (exerciseDir, testsDir string) {
	t.Helper()
	root := t.TempDir()
	exerciseDir = filepath.Join(root, "exercises", "double-01-go")
	testsDir = filepath.Join(root, "tests", "double-01-go")

	writeFile(t, filepath.Join(exerciseDir, "exercise.json"), `{
  "id": "double-01-go",
  "problem_id": "double-01",
  "title": "Double",
  "category": "dsa",
  "language": "go",
  "time_limit_min": 10,
  "tutor_mode": "hints-first",
  "repo_path": "./repo",
  "test_command": "go test ./..."
}`)
	writeFile(t, filepath.Join(exerciseDir, "repo", "go.mod"), goModContents)
	writeFile(t, filepath.Join(exerciseDir, "repo", "solution.go"), starterBody)
	writeFile(t, filepath.Join(testsDir, "solution_test.go"), goTestContents)

	if referenceBody != "" {
		writeFile(t, filepath.Join(exerciseDir, ".reference", "solution.go"), referenceBody)
	}
	return exerciseDir, testsDir
}

const correctDouble = "package main\n\nfunc Double(n int) int { return n * 2 }\n"
const blankDouble = "package main\n\nfunc Double(n int) int { return 0 }\n"
const wrongReferenceDouble = "package main\n\nfunc Double(n int) int { return n + 1 }\n"
const accidentallyCorrectStarter = "package main\n\nfunc Double(n int) int { return n * 2 }\n"

// fixtureDesignExercise builds a design-kind exercise: exercise.json,
// repo/{problem.md, solution.md}, and tests/<id>/rubric.md -- each
// optional so tests can knock out one structural piece at a time.
func fixtureDesignExercise(t *testing.T, withProblem, withSolution, withRubric bool) (exerciseDir, testsDir string) {
	t.Helper()
	root := t.TempDir()
	exerciseDir = filepath.Join(root, "exercises", "url-shortener-01-coach")
	testsDir = filepath.Join(root, "tests", "url-shortener-01-coach")

	writeFile(t, filepath.Join(exerciseDir, "exercise.json"), `{
  "id": "url-shortener-01-coach",
  "problem_id": "url-shortener-01",
  "title": "Design Pastebin / Bit.ly",
  "kind": "design",
  "category": "system-design",
  "language": "coach",
  "time_limit_min": 90,
  "tutor_mode": "design-coach",
  "repo_path": "./repo",
  "test_command": ""
}`)
	if withProblem {
		writeFile(t, filepath.Join(exerciseDir, "repo", "problem.md"), "# Design Pastebin\n\nprompt")
	}
	if withSolution {
		writeFile(t, filepath.Join(exerciseDir, "repo", "solution.md"), "# My design\n\n## Step 1")
	}
	if !withProblem && !withSolution {
		// repo_path must exist for exercise.Load either way.
		if err := os.MkdirAll(filepath.Join(exerciseDir, "repo"), 0o755); err != nil {
			t.Fatalf("mkdir repo: %v", err)
		}
	}
	if withRubric {
		writeFile(t, filepath.Join(testsDir, "rubric.md"), "- estimates\n- high-level design")
	}
	return exerciseDir, testsDir
}

func TestExercise_DesignWithAllStructuralPieces_IsOK(t *testing.T) {
	exDir, testsDir := fixtureDesignExercise(t, true, true, true)

	result, err := Exercise(exDir, testsDir)
	if err != nil {
		t.Fatalf("Exercise: %v", err)
	}
	if !result.OK() {
		t.Error("expected OK() for a structurally complete design exercise")
	}
}

func TestExercise_DesignMissingPieces_EachFailsLoud(t *testing.T) {
	cases := []struct {
		name                                  string
		withProblem, withSolution, withRubric bool
		wantErr                               string
	}{
		{"missing problem.md", false, true, true, "problem.md"},
		{"missing solution.md template", true, false, true, "solution.md"},
		{"missing rubric.md", true, true, false, "rubric.md"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			exDir, testsDir := fixtureDesignExercise(t, tc.withProblem, tc.withSolution, tc.withRubric)
			_, err := Exercise(exDir, testsDir)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Exercise err = %v, want error naming %q -- a missing rubric would otherwise only surface as a live 30s submit timeout", err, tc.wantErr)
			}
		})
	}
}

func TestExercise_ReferencePassesAndStarterFails_IsOK(t *testing.T) {
	exDir, testsDir := fixtureExercise(t, blankDouble, correctDouble)

	result, err := Exercise(exDir, testsDir)
	if err != nil {
		t.Fatalf("Exercise: %v", err)
	}
	if !result.StarterFailed {
		t.Errorf("expected the blank starter to fail the hidden tests, output:\n%s", result.StarterOutput)
	}
	if !result.ReferencePassed {
		t.Errorf("expected the reference solution to pass the hidden tests, output:\n%s", result.ReferenceOutput)
	}
	if !result.OK() {
		t.Error("expected OK() to be true when starter fails and reference passes")
	}
}

func TestExercise_StarterAccidentallyPassing_IsCaught(t *testing.T) {
	// The whole point of this tool: if the starter already satisfies the
	// hidden tests (a vacuous or too-weak test suite), that must be
	// flagged, not silently accepted.
	exDir, testsDir := fixtureExercise(t, accidentallyCorrectStarter, correctDouble)

	result, err := Exercise(exDir, testsDir)
	if err != nil {
		t.Fatalf("Exercise: %v", err)
	}
	if result.StarterFailed {
		t.Error("expected StarterFailed=false — the starter unexpectedly already passes")
	}
	if result.OK() {
		t.Error("expected OK()=false when the starter unexpectedly passes")
	}
}

func TestExercise_BrokenReference_IsCaught(t *testing.T) {
	exDir, testsDir := fixtureExercise(t, blankDouble, wrongReferenceDouble)

	result, err := Exercise(exDir, testsDir)
	if err != nil {
		t.Fatalf("Exercise: %v", err)
	}
	if result.ReferencePassed {
		t.Error("expected ReferencePassed=false — the reference solution is actually wrong")
	}
	if result.OK() {
		t.Error("expected OK()=false when the reference solution fails")
	}
	if !strings.Contains(result.ReferenceOutput, "FAIL") {
		t.Errorf("expected go test's own FAIL output to be captured, got:\n%s", result.ReferenceOutput)
	}
}

func TestExercise_MissingReferenceDir_ReturnsError(t *testing.T) {
	exDir, testsDir := fixtureExercise(t, blankDouble, "" /* no .reference */)

	_, err := Exercise(exDir, testsDir)
	if err == nil {
		t.Fatal("expected an error when .reference is missing, got nil")
	}
}

func TestExercise_InvalidExerciseJSON_ReturnsError(t *testing.T) {
	root := t.TempDir()
	exDir := filepath.Join(root, "exercises", "broken-01")
	writeFile(t, filepath.Join(exDir, "exercise.json"), `not json`)

	_, err := Exercise(exDir, filepath.Join(root, "tests", "broken-01"))
	if err == nil {
		t.Fatal("expected an error for invalid exercise.json, got nil")
	}
}

func TestExercise_DoesNotLeakWorkspacesOnDisk(t *testing.T) {
	// buildWorkspace uses os.MkdirTemp — Exercise must clean both
	// temp workspaces it creates (starter run + reference run) rather
	// than leaving them behind on every verify pass across 150+ exercises.
	before, err := filepath.Glob(filepath.Join(os.TempDir(), "verify-workspace-*"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}

	exDir, testsDir := fixtureExercise(t, blankDouble, correctDouble)
	if _, err := Exercise(exDir, testsDir); err != nil {
		t.Fatalf("Exercise: %v", err)
	}

	after, err := filepath.Glob(filepath.Join(os.TempDir(), "verify-workspace-*"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}
	if len(after) != len(before) {
		t.Errorf("expected no leaked verify-workspace-* dirs, before=%d after=%d", len(before), len(after))
	}
}
