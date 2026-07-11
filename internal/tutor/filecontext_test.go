package tutor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildFileContext_ReturnsExactContents(t *testing.T) {
	dir := t.TempDir()
	want := "package main\n\nfunc solve() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte(want), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	got := buildFileContext(dir, 8000)
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
	if !strings.HasPrefix(got, strings.Repeat("x", 10)) {
		t.Errorf("buildFileContext = %q, want it to start with the first 10 bytes", got)
	}
	if !strings.Contains(got, "truncated") {
		t.Errorf("buildFileContext = %q, want a truncation marker", got)
	}
	if !strings.Contains(got, "100 bytes") {
		t.Errorf("buildFileContext = %q, want the marker to mention the real file size (100 bytes)", got)
	}
	// The raw content prefix must never exceed maxBytes, even though the
	// marker text appended after it does.
	rawPrefix := strings.SplitN(got, "\n...[truncated", 2)[0]
	if len(rawPrefix) > 10 {
		t.Errorf("raw content prefix is %d bytes, want <= 10", len(rawPrefix))
	}
}

func TestBuildFileContext_RereadsFreshEachCall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "solution.cpp")
	if err := os.WriteFile(path, []byte("first version"), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}

	first := buildFileContext(dir, 8000)
	if first != "first version" {
		t.Fatalf("first read = %q, want %q", first, "first version")
	}

	if err := os.WriteFile(path, []byte("second version"), 0o644); err != nil {
		t.Fatalf("rewrite solution file: %v", err)
	}

	second := buildFileContext(dir, 8000)
	if second != "second version" {
		t.Errorf("second read = %q, want %q (should re-read, not cache)", second, "second version")
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
