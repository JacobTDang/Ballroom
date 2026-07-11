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

func TestHighlightLinesTool_SucceedsAgainstLiveNvim(t *testing.T) {
	socket := startTestNvim(t)

	tl, err := newHighlightLinesTool(Config{NvimSocket: socket})
	if err != nil {
		t.Fatalf("newHighlightLinesTool: %v", err)
	}

	in, err := json.Marshal(highlightLinesInput{File: "test.txt", Start: 1, End: 1, Note: "off by one here"})
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}
	out, err := tl.InvokableRun(context.Background(), string(in))
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result highlightLinesOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Status != "ok" {
		t.Errorf("Status = %q, want %q", result.Status, "ok")
	}
}

func TestHighlightLinesTool_NoEditorReturnsFriendlyStatus(t *testing.T) {
	tl, err := newHighlightLinesTool(Config{NvimSocket: ""})
	if err != nil {
		t.Fatalf("newHighlightLinesTool: %v", err)
	}

	in, err := json.Marshal(highlightLinesInput{File: "test.txt", Start: 1, End: 1, Note: "a note"})
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}
	out, err := tl.InvokableRun(context.Background(), string(in))
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result highlightLinesOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Status == "ok" || result.Status == "" {
		t.Errorf("Status = %q, want a friendly note that no editor is attached", result.Status)
	}
}

func TestHighlightLinesTool_OutOfBoundsRangeReturnsError(t *testing.T) {
	socket := startTestNvim(t)

	tl, err := newHighlightLinesTool(Config{NvimSocket: socket})
	if err != nil {
		t.Fatalf("newHighlightLinesTool: %v", err)
	}

	// A freshly headless-started nvim's scratch buffer has exactly 1
	// line, so a range starting at line 999 is out of bounds and
	// ballroom_highlight.lua's add_highlight must report an error.
	in, err := json.Marshal(highlightLinesInput{File: "test.txt", Start: 999, End: 999, Note: "n/a"})
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}
	if _, err := tl.InvokableRun(context.Background(), string(in)); err == nil {
		t.Error("expected an error for an out-of-bounds line range, got nil")
	}
}

func TestReadCursorPositionTool_ReturnsPositionAgainstLiveNvim(t *testing.T) {
	socket := startTestNvim(t)

	tl, err := newReadCursorPositionTool(Config{NvimSocket: socket})
	if err != nil {
		t.Fatalf("newReadCursorPositionTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result cursorPositionOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if !result.Available {
		t.Fatal("Available = false, want true against a live nvim")
	}
	if result.Line != 1 {
		t.Errorf("Line = %d, want 1 (freshly started nvim's scratch buffer)", result.Line)
	}
}

func TestReadCursorPositionTool_NoEditorReturnsUnavailable(t *testing.T) {
	tl, err := newReadCursorPositionTool(Config{NvimSocket: ""})
	if err != nil {
		t.Fatalf("newReadCursorPositionTool: %v", err)
	}

	out, err := tl.InvokableRun(context.Background(), "{}")
	if err != nil {
		t.Fatalf("InvokableRun: %v", err)
	}

	var result cursorPositionOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unmarshal tool output %q: %v", out, err)
	}
	if result.Available {
		t.Error("Available = true, want false when no editor is attached")
	}
}

func TestBuildTools_ReturnsAllFiveTools(t *testing.T) {
	tools, err := buildTools(Config{WorkDir: t.TempDir(), MaxContextBytes: 8000})
	if err != nil {
		t.Fatalf("buildTools: %v", err)
	}
	if len(tools) != 5 {
		t.Fatalf("buildTools returned %d tools, want 5", len(tools))
	}

	names := make(map[string]bool)
	for _, tl := range tools {
		info, err := tl.Info(context.Background())
		if err != nil {
			t.Fatalf("tool.Info: %v", err)
		}
		names[info.Name] = true
	}
	for _, want := range []string{
		"read_solution_file", "read_problem_statement", "read_test_output",
		"highlight_lines", "read_cursor_position",
	} {
		if !names[want] {
			t.Errorf("buildTools missing tool %q, got names %v", want, names)
		}
	}
}
