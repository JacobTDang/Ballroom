package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

// seedExercise builds a real, runnable go-based exercise under root
// (exercises/<id>/ + tests/<id>/), with the given starter/reference
// bodies (referenceBody == "" omits .reference entirely).
func seedExercise(t *testing.T, root, id, starterBody, referenceBody string) {
	t.Helper()
	exDir := filepath.Join(root, "exercises", id)
	testsDir := filepath.Join(root, "tests", id)

	writeFile(t, filepath.Join(exDir, "exercise.json"), `{
  "id": "`+id+`",
  "title": "Double",
  "category": "dsa",
  "language": "go",
  "time_limit_min": 10,
  "tutor_mode": "hints-first",
  "repo_path": "./repo",
  "test_command": "go test ./..."
}`)
	writeFile(t, filepath.Join(exDir, "repo", "go.mod"), goModContents)
	writeFile(t, filepath.Join(exDir, "repo", "solution.go"), starterBody)
	writeFile(t, filepath.Join(testsDir, "solution_test.go"), goTestContents)
	if referenceBody != "" {
		writeFile(t, filepath.Join(exDir, ".reference", "solution.go"), referenceBody)
	}
}

const correctDouble = "package main\n\nfunc Double(n int) int { return n * 2 }\n"
const blankDouble = "package main\n\nfunc Double(n int) int { return 0 }\n"

func TestListExerciseDirs_SkipsTemplateAndNonExerciseDirs(t *testing.T) {
	root := t.TempDir()
	seedExercise(t, root, "double-01-go", blankDouble, correctDouble)
	writeFile(t, filepath.Join(root, "exercises", "_template", "exercise.json"), "{}")
	if err := os.MkdirAll(filepath.Join(root, "exercises", "not-an-exercise"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	ids, err := listExerciseDirs(filepath.Join(root, "exercises"))
	if err != nil {
		t.Fatalf("listExerciseDirs: %v", err)
	}
	if len(ids) != 1 || ids[0] != "double-01-go" {
		t.Errorf("ids = %v, want [double-01-go]", ids)
	}
}

func TestRun_AllPass_ReturnsOKTrue(t *testing.T) {
	root := t.TempDir()
	seedExercise(t, root, "double-01-go", blankDouble, correctDouble)

	var out bytes.Buffer
	ok, err := run(root, &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !ok {
		t.Errorf("expected ok=true, output:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "PASS double-01-go") {
		t.Errorf("expected a PASS line for double-01-go, got:\n%s", out.String())
	}
}

func TestRun_MissingReference_ReturnsOKFalseAndReportsIt(t *testing.T) {
	root := t.TempDir()
	seedExercise(t, root, "double-01-go", blankDouble, "" /* no .reference */)

	var out bytes.Buffer
	ok, err := run(root, &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if ok {
		t.Error("expected ok=false when an exercise has no .reference")
	}
	if !strings.Contains(out.String(), "double-01-go") {
		t.Errorf("expected the failing exercise id to be reported, got:\n%s", out.String())
	}
}

func TestRun_StarterAccidentallyPasses_ReturnsOKFalse(t *testing.T) {
	root := t.TempDir()
	seedExercise(t, root, "double-01-go", correctDouble /* starter == reference */, correctDouble)

	var out bytes.Buffer
	ok, err := run(root, &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if ok {
		t.Error("expected ok=false when the starter unexpectedly already passes")
	}
}
