// Package exercise loads exercise definitions from exercises/<id>/exercise.json.
// Hidden tests are never read from here — they live under the sibling
// tests/<id>/ tree and are only touched by the orchestrator at submit time.
package exercise

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CategoryPattern        = "pattern"
	CategoryDebug          = "debug"
	CategoryConcurrency    = "concurrency"
	CategoryImplementation = "implementation"
	CategoryAIAssisted     = "ai-assisted"
)

const (
	LanguageGo     = "go"
	LanguageCpp    = "cpp"
	LanguagePython = "python"
)

const (
	TutorModeSyntaxOnly = "syntax-only"
	TutorModeHintsFirst = "hints-first"
	TutorModeFullAssist = "full-assist"
)

var validCategories = map[string]bool{
	CategoryPattern:        true,
	CategoryDebug:          true,
	CategoryConcurrency:    true,
	CategoryImplementation: true,
	CategoryAIAssisted:     true,
}

var validLanguages = map[string]bool{
	LanguageGo:     true,
	LanguageCpp:    true,
	LanguagePython: true,
}

var validTutorModes = map[string]bool{
	TutorModeSyntaxOnly: true,
	TutorModeHintsFirst: true,
	TutorModeFullAssist: true,
}

// Exercise is a parsed, validated exercise.json.
type Exercise struct {
	ID           string
	Title        string
	Category     string
	Language     string
	TimeLimitMin int
	TutorMode    string
	RepoPath     string // resolved to an absolute path
	TestCommand  string
}

// raw mirrors the on-disk JSON shape before validation/path resolution.
type raw struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	Language     string `json:"language"`
	TimeLimitMin int    `json:"time_limit_min"`
	TutorMode    string `json:"tutor_mode"`
	RepoPath     string `json:"repo_path"`
	TestCommand  string `json:"test_command"`
}

// Load reads and validates the exercise definition at path (exercise.json).
// repo_path is resolved relative to path's directory.
func Load(path string) (Exercise, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Exercise{}, fmt.Errorf("exercise: read %s: %w", path, err)
	}

	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return Exercise{}, fmt.Errorf("exercise: parse %s: %w", path, err)
	}

	if r.ID == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: id is required", path)
	}
	if r.Title == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: title is required", path)
	}
	if !validCategories[r.Category] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid category %q", path, r.Category)
	}
	if !validLanguages[r.Language] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid language %q", path, r.Language)
	}
	if r.TimeLimitMin <= 0 {
		return Exercise{}, fmt.Errorf("exercise: %s: time_limit_min must be > 0", path)
	}
	if !validTutorModes[r.TutorMode] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid tutor_mode %q", path, r.TutorMode)
	}
	if r.RepoPath == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: repo_path is required", path)
	}
	if r.TestCommand == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: test_command is required", path)
	}

	repoPath := r.RepoPath
	if !filepath.IsAbs(repoPath) {
		repoPath = filepath.Join(filepath.Dir(path), repoPath)
	}
	if info, err := os.Stat(repoPath); err != nil || !info.IsDir() {
		return Exercise{}, fmt.Errorf("exercise: %s: repo_path %q does not exist or is not a directory", path, repoPath)
	}

	return Exercise{
		ID:           r.ID,
		Title:        r.Title,
		Category:     r.Category,
		Language:     r.Language,
		TimeLimitMin: r.TimeLimitMin,
		TutorMode:    r.TutorMode,
		RepoPath:     repoPath,
		TestCommand:  r.TestCommand,
	}, nil
}
