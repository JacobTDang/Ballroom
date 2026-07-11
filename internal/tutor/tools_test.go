package tutor

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadSolutionFileTool_ReturnsFileContents(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "solution.py"), []byte("def solve(): pass"), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	tl, err := newReadSolutionFileTool(Config{WorkDir: dir, MaxContextBytes: 8000})
	if err != nil {
		t.Fatalf("newReadSolutionFileTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result readFileOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Content != "def solve(): pass" {
		t.Errorf("Content = %q, want %q", result.Content, "def solve(): pass")
	}
}

func TestReadSolutionFileTool_NoFileReturnsFriendlyNote(t *testing.T) {
	dir := t.TempDir()

	tl, err := newReadSolutionFileTool(Config{WorkDir: dir, MaxContextBytes: 8000})
	if err != nil {
		t.Fatalf("newReadSolutionFileTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result readFileOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Content == "" {
		t.Error("Content is empty, want a friendly note explaining there's no solution file yet")
	}
}

func TestReadProblemStatementTool_ReturnsProblemMd(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "problem.md"), []byte("# Reverse a String"), 0o644); err != nil {
		t.Fatalf("write problem.md: %v", err)
	}

	tl, err := newReadProblemStatementTool(Config{WorkDir: dir})
	if err != nil {
		t.Fatalf("newReadProblemStatementTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result readFileOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Content != "# Reverse a String" {
		t.Errorf("Content = %q, want %q", result.Content, "# Reverse a String")
	}
}

func TestReadTestOutputTool_ReturnsResultWhenPresent(t *testing.T) {
	dir := t.TempDir()
	data, err := json.Marshal(lastTestResult{
		Result:      "fail",
		Output:      "FAIL: something broke",
		TestCommand: "python3 -m pytest -q",
	})
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, lastTestResultFile), data, 0o644); err != nil {
		t.Fatalf("write last test result file: %v", err)
	}

	tl, err := newReadTestOutputTool(Config{WorkDir: dir})
	if err != nil {
		t.Fatalf("newReadTestOutputTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result readTestOutputOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if !result.Available {
		t.Error("Available = false, want true")
	}
	if result.Result != "fail" {
		t.Errorf("Result = %q, want %q", result.Result, "fail")
	}
	if result.Output != "FAIL: something broke" {
		t.Errorf("Output = %q, want %q", result.Output, "FAIL: something broke")
	}
}

func TestReadTestOutputTool_NoResultReturnsUnavailable(t *testing.T) {
	dir := t.TempDir() // never submitted

	tl, err := newReadTestOutputTool(Config{WorkDir: dir})
	if err != nil {
		t.Fatalf("newReadTestOutputTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result readTestOutputOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Available {
		t.Error("Available = true, want false when no submission has happened yet")
	}
}

func TestBuildTools_ReturnsAllReadOnlyTools(t *testing.T) {
	tools, err := buildTools(Config{WorkDir: t.TempDir(), MaxContextBytes: 8000})
	if err != nil {
		t.Fatalf("buildTools: %v", err)
	}
	if len(tools) != 3 {
		t.Fatalf("buildTools returned %d tools, want 3", len(tools))
	}

	names := make(map[string]bool)
	for _, tl := range tools {
		info, err := tl.Info(context.Background())
		if err != nil {
			t.Fatalf("tool.Info: %v", err)
		}
		names[info.Name] = true
	}
	for _, want := range []string{"read_solution_file", "read_problem_statement", "read_test_output"} {
		if !names[want] {
			t.Errorf("buildTools missing tool %q, got names %v", want, names)
		}
	}
}
