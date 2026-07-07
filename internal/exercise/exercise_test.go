package exercise

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeExercise(t *testing.T, dir string, fields map[string]any) string {
	t.Helper()
	base := map[string]any{
		"id":             "two-pointers-01",
		"title":          "Two Sum II",
		"category":       "pattern",
		"language":       "go",
		"time_limit_min": 25,
		"tutor_mode":     "hints-first",
		"repo_path":      "./repo",
		"test_command":   "go test ./...",
	}
	for k, v := range fields {
		base[k] = v
	}
	data, err := json.Marshal(base)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	path := filepath.Join(dir, "exercise.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestLoad_ValidExercise(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	path := writeExercise(t, dir, nil)

	ex, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if ex.ID != "two-pointers-01" {
		t.Errorf("ID = %q, want %q", ex.ID, "two-pointers-01")
	}
	if ex.Category != CategoryPattern {
		t.Errorf("Category = %q, want %q", ex.Category, CategoryPattern)
	}
	if ex.Language != LanguageGo {
		t.Errorf("Language = %q, want %q", ex.Language, LanguageGo)
	}
	if ex.TimeLimitMin != 25 {
		t.Errorf("TimeLimitMin = %d, want 25", ex.TimeLimitMin)
	}
	if ex.TutorMode != TutorModeHintsFirst {
		t.Errorf("TutorMode = %q, want %q", ex.TutorMode, TutorModeHintsFirst)
	}
	if ex.TestCommand != "go test ./..." {
		t.Errorf("TestCommand = %q, want %q", ex.TestCommand, "go test ./...")
	}

	wantRepo := filepath.Join(dir, "repo")
	if ex.RepoPath != wantRepo {
		t.Errorf("RepoPath = %q, want %q (resolved relative to exercise.json's dir)", ex.RepoPath, wantRepo)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "does-not-exist.json")); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidCategory(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"category": "not-a-real-category"})

	if _, err := Load(path); err == nil {
		t.Fatal("expected error for invalid category, got nil")
	}
}

func TestLoad_InvalidLanguage(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"language": "rust"})

	if _, err := Load(path); err == nil {
		t.Fatal("expected error for invalid language, got nil")
	}
}

func TestLoad_InvalidTutorMode(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"tutor_mode": "godmode"})

	if _, err := Load(path); err == nil {
		t.Fatal("expected error for invalid tutor_mode, got nil")
	}
}

func TestLoad_ZeroTimeLimit(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"time_limit_min": 0})

	if _, err := Load(path); err == nil {
		t.Fatal("expected error for zero time_limit_min, got nil")
	}
}

func TestLoad_MissingRepoDir(t *testing.T) {
	dir := t.TempDir()
	// deliberately not creating dir/repo
	path := writeExercise(t, dir, nil)

	if _, err := Load(path); err == nil {
		t.Fatal("expected error when repo_path does not exist, got nil")
	}
}

func TestLoad_EmptyTestCommand(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"test_command": ""})

	if _, err := Load(path); err == nil {
		t.Fatal("expected error for empty test_command, got nil")
	}
}
