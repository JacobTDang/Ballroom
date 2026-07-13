package tutor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBuildFileContext_NumbersEachLine locks in the fix for highlight_lines
// landing on the wrong line: without ground-truth line numbers in what the
// model reads back, it has to count lines itself from raw prose-formatted
// code (blank lines, comments, wrapped display width, ...), which is
// exactly the kind of thing a model gets wrong. Numbering matches
// highlight_lines' own 1-indexed start/end contract directly, so the
// model can copy a number instead of counting one.
func TestBuildFileContext_NumbersEachLine(t *testing.T) {
	dir := t.TempDir()
	content := "package main\n\nfunc solve() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	got := buildFileContext(dir, 8000)
	want := "1\tpackage main\n2\t\n3\tfunc solve() {}"
	if got != want {
		t.Errorf("buildFileContext = %q, want %q", got, want)
	}
}

// TestBuildFileContext_NumbersLastLineWithNoTrailingNewline covers a file
// that doesn't end in a newline (e.g. mid-edit) -- the last line must
// still get a number, not get silently dropped or merged into the
// previous one.
func TestBuildFileContext_NumbersLastLineWithNoTrailingNewline(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte("a\nb\nc"), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	got := buildFileContext(dir, 8000)
	want := "1\ta\n2\tb\n3\tc"
	if got != want {
		t.Errorf("buildFileContext = %q, want %q", got, want)
	}
}

func TestBuildFileContext_MissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir() // no solution.* file at all

	got := buildFileContext(dir, 8000)
	if got != "" {
		t.Errorf("buildFileContext = %q, want empty string for a missing solution file", got)
	}
}

func TestBuildFileContext_TruncatesOversizedFile(t *testing.T) {
	dir := t.TempDir()
	big := strings.Repeat("x", 100)
	if err := os.WriteFile(filepath.Join(dir, "solution.py"), []byte(big), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	got := buildFileContext(dir, 10)
	wantPrefix := "1\t" + strings.Repeat("x", 10)
	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("buildFileContext = %q, want it to start with %q", got, wantPrefix)
	}
	if !strings.Contains(got, "truncated") {
		t.Errorf("buildFileContext = %q, want a truncation marker", got)
	}
	if !strings.Contains(got, "100 bytes") {
		t.Errorf("buildFileContext = %q, want the marker to mention the real file size (100 bytes)", got)
	}
	// The numbered content before the marker must be exactly the
	// (numbered) first 10 raw bytes, even though the marker text appended
	// after it isn't itself bounded by maxBytes.
	rawPrefix := strings.SplitN(got, "\n...[truncated", 2)[0]
	if rawPrefix != wantPrefix {
		t.Errorf("numbered content before the truncation marker = %q, want exactly %q", rawPrefix, wantPrefix)
	}
}

func TestBuildFileContext_RereadsFreshEachCall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "solution.cpp")
	if err := os.WriteFile(path, []byte("first version"), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	first := buildFileContext(dir, 8000)
	if first != "1\tfirst version" {
		t.Fatalf("first read = %q, want %q", first, "1\tfirst version")
	}

	if err := os.WriteFile(path, []byte("second version"), 0o644); err != nil {
		t.Fatalf("rewrite solution file: %v", err)
	}

	second := buildFileContext(dir, 8000)
	if second != "1\tsecond version" {
		t.Errorf("second read = %q, want %q (should re-read, not cache)", second, "1\tsecond version")
	}
}

func TestReadProblemStatement_ReturnsContents(t *testing.T) {
	dir := t.TempDir()
	want := "# Two Sum\n\nGiven an array...\n"
	if err := os.WriteFile(filepath.Join(dir, "problem.md"), []byte(want), 0o644); err != nil {
		t.Fatalf("write problem.md: %v", err)
	}

	got := readProblemStatement(dir)
	if got != want {
		t.Errorf("readProblemStatement = %q, want %q", got, want)
	}
}

func TestReadProblemStatement_MissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir() // no problem.md

	got := readProblemStatement(dir)
	if got != "" {
		t.Errorf("readProblemStatement = %q, want empty string when problem.md is missing", got)
	}
}

func TestReadLastTestResult_ReturnsWrittenResult(t *testing.T) {
	dir := t.TempDir()
	recordedAt := time.Date(2026, 7, 10, 18, 4, 11, 0, time.UTC)
	data, err := json.Marshal(lastTestResult{
		Result:      "pass",
		Output:      "ok\nPASS",
		TestCommand: "go test ./...",
		RecordedAt:  recordedAt,
	})
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, lastTestResultFile), data, 0o644); err != nil {
		t.Fatalf("write last test result file: %v", err)
	}

	got, ok, err := readLastTestResult(dir)
	if err != nil {
		t.Fatalf("readLastTestResult: %v", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if got.Result != "pass" {
		t.Errorf("Result = %q, want %q", got.Result, "pass")
	}
	if got.Output != "ok\nPASS" {
		t.Errorf("Output = %q, want %q", got.Output, "ok\nPASS")
	}
	if got.TestCommand != "go test ./..." {
		t.Errorf("TestCommand = %q, want %q", got.TestCommand, "go test ./...")
	}
	if !got.RecordedAt.Equal(recordedAt) {
		t.Errorf("RecordedAt = %v, want %v", got.RecordedAt, recordedAt)
	}
}

func TestReadLastTestResult_MissingFileReturnsNotAvailable(t *testing.T) {
	dir := t.TempDir() // never submitted, or sandbox mode

	got, ok, err := readLastTestResult(dir)
	if err != nil {
		t.Fatalf("readLastTestResult: %v", err)
	}
	if ok {
		t.Errorf("ok = true, want false for a missing file (got %+v)", got)
	}
}

func TestReadLastTestResult_MalformedFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, lastTestResultFile), []byte("not json"), 0o644); err != nil {
		t.Fatalf("write malformed file: %v", err)
	}

	_, ok, err := readLastTestResult(dir)
	if err == nil {
		t.Fatal("expected an error for a malformed (but present) result file, got nil")
	}
	if ok {
		t.Error("ok = true, want false alongside the error")
	}
}
