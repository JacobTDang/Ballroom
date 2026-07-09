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
	CategoryDSA            = "dsa"
	CategoryDebug          = "debug"
	CategoryConcurrency    = "concurrency"
	CategoryImplementation = "implementation"
	CategoryAIAssisted     = "ai-assisted"
)

// The NeetCode 150 roadmap's own category breakdown — every problem
// generated from that list lands in one of these, not the generic
// CategoryDSA bucket (see internal/catalog.DisplayCategory for how each
// renders).
const (
	CategoryArraysHashing   = "arrays-hashing"
	CategoryTwoPointers     = "two-pointers"
	CategorySlidingWindow   = "sliding-window"
	CategoryStack           = "stack"
	CategoryBinarySearch    = "binary-search"
	CategoryLinkedList      = "linked-list"
	CategoryTrees           = "trees"
	CategoryTries           = "tries"
	CategoryHeap            = "heap"
	CategoryBacktracking    = "backtracking"
	CategoryGraphs          = "graphs"
	CategoryAdvancedGraphs  = "advanced-graphs"
	CategoryDP1D            = "1d-dp"
	CategoryDP2D            = "2d-dp"
	CategoryGreedy          = "greedy"
	CategoryIntervals       = "intervals"
	CategoryMathGeometry    = "math-geometry"
	CategoryBitManipulation = "bit-manipulation"
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
	CategoryDSA:            true,
	CategoryDebug:          true,
	CategoryConcurrency:    true,
	CategoryImplementation: true,
	CategoryAIAssisted:     true,

	CategoryArraysHashing:   true,
	CategoryTwoPointers:     true,
	CategorySlidingWindow:   true,
	CategoryStack:           true,
	CategoryBinarySearch:    true,
	CategoryLinkedList:      true,
	CategoryTrees:           true,
	CategoryTries:           true,
	CategoryHeap:            true,
	CategoryBacktracking:    true,
	CategoryGraphs:          true,
	CategoryAdvancedGraphs:  true,
	CategoryDP1D:            true,
	CategoryDP2D:            true,
	CategoryGreedy:          true,
	CategoryIntervals:       true,
	CategoryMathGeometry:    true,
	CategoryBitManipulation: true,
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
	ProblemID    string // groups language variants of the same problem
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
	ProblemID    string `json:"problem_id"`
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

	problemID := r.ProblemID
	if problemID == "" {
		// Standalone exercises with no language siblings don't need to
		// declare a problem_id just to load — they're their own problem.
		problemID = r.ID
	}

	return Exercise{
		ID:           r.ID,
		ProblemID:    problemID,
		Title:        r.Title,
		Category:     r.Category,
		Language:     r.Language,
		TimeLimitMin: r.TimeLimitMin,
		TutorMode:    r.TutorMode,
		RepoPath:     repoPath,
		TestCommand:  r.TestCommand,
	}, nil
}
