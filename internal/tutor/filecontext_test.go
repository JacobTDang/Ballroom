package tutor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
