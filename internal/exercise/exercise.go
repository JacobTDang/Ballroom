// Package exercise loads exercise definitions from exercises/<id>/exercise.json.
// Hidden tests are never read from here — they live under the sibling
// tests/<id>/ tree and are only touched by the orchestrator at submit time.
package exercise

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// Kind separates the two session shapes an exercise can describe:
// coding (the default -- solution code, hidden tests, pass/fail runs)
// and design (a system-design question -- the "solution" is a design
// doc, the hidden tests/<id>/ tree holds a grading rubric instead of
// tests, and submit is self-assessed). Empty kind in exercise.json
// means coding, so none of the pre-existing exercise files change.
const (
	KindCoding = "coding"
	KindDesign = "design"
)

// CategorySystemDesign holds the design-kind questions (see
// docs/system-design-roadmap.md). Deliberately NOT a DSA subcategory --
// it's its own top-level practice-picker entry.
const CategorySystemDesign = "system-design"

// CategoryOODesign holds the object-oriented design problems from the
// same roadmap (hash map, parking lot, call center, ...). These are
// ordinary coding exercises with real hidden tests -- KindCoding, not
// KindDesign -- grouped as their own top-level category.
const CategoryOODesign = "oo-design"

// CategoryBehavioral holds the behavioral-interview questions
// ("tell me about a time..."): design-kind sessions whose solution.md
// is a STAR story, graded against a STAR rubric on submit.
const CategoryBehavioral = "behavioral"

const (
	LanguageGo     = "go"
	LanguageCpp    = "cpp"
	LanguagePython = "python"
)

// Design exercises have no programming language; the language slot
// instead carries the session style, which is what actually varies
// between sibling variants of the same design question. This reuses the
// ProblemID variant grouping and the TUI's language picker unchanged --
// picking "coach" vs "interviewer" IS picking the exercise variant.
const (
	LanguageCoach       = "coach"
	LanguageInterviewer = "interviewer"
)

const (
	TutorModeSyntaxOnly = "syntax-only"
	TutorModeHintsFirst = "hints-first"
	TutorModeFullAssist = "full-assist"
)

// Design-session tutor personas (internal/tutor/prompts.go): the
// interviewer runs a probing mock interview and never restates the
// problem (that's the candidate's job); the design coach walks through
// the 4-step method with the user. Only valid on design-kind exercises.
const (
	TutorModeInterviewer = "interviewer"
	TutorModeDesignCoach = "design-coach"
)

// Behavioral-session tutor personas: the same design-kind machinery
// (solution.md, rubric reveal, LLM grading) with STAR stories instead
// of architectures -- the behavioral interviewer asks the question and
// probes for specifics, the story coach builds the story one STAR
// section at a time. Only valid on design-kind exercises.
const (
	TutorModeBehavioralInterviewer = "behavioral-interviewer"
	TutorModeStoryCoach            = "story-coach"
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

	CategorySystemDesign: true,
	CategoryOODesign:     true,
	CategoryBehavioral:   true,
}

// The language and tutor-mode vocabularies are kind-gated: a coding
// exercise can't claim a session style, and a design exercise can't
// claim a programming language or a coding tutor mode. Load picks the
// map by kind.
var validLanguages = map[string]bool{
	LanguageGo:     true,
	LanguageCpp:    true,
	LanguagePython: true,
}

var validDesignLanguages = map[string]bool{
	LanguageCoach:       true,
	LanguageInterviewer: true,
}

var validTutorModes = map[string]bool{
	TutorModeSyntaxOnly: true,
	TutorModeHintsFirst: true,
	TutorModeFullAssist: true,
}

var validDesignTutorModes = map[string]bool{
	TutorModeInterviewer:           true,
	TutorModeDesignCoach:           true,
	TutorModeBehavioralInterviewer: true,
	TutorModeStoryCoach:            true,
}

// Exercise is a parsed, validated exercise.json.
type Exercise struct {
	ID           string
	ProblemID    string // groups language variants of the same problem
	Title        string
	Kind         string // KindCoding (default) or KindDesign
	Category     string
	Language     string // programming language, or session style for design kind
	TimeLimitMin int
	TutorMode    string
	RepoPath     string // resolved to an absolute path
	TestCommand  string // empty for design kind (nothing to run)
	// VideoURL optionally links a solution walkthrough video, shown in
	// the problem statement footer and with submit results. Empty means
	// none.
	VideoURL string
	// Difficulty optionally rates the problem — one of DifficultyEasy/
	// Medium/Hard, or empty for unrated. Lowercase-only on disk: the
	// NeetCode site data's "Easy"/"Medium"/"Hard" is normalized at fill
	// time so every consumer (picker badges, sorting) matches one
	// vocabulary.
	Difficulty string
}

// The difficulty vocabulary. Values are the on-disk strings.
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)

var validDifficulties = map[string]bool{
	DifficultyEasy:   true,
	DifficultyMedium: true,
	DifficultyHard:   true,
}

// raw mirrors the on-disk JSON shape before validation/path resolution.
type raw struct {
	ID           string `json:"id"`
	ProblemID    string `json:"problem_id"`
	Title        string `json:"title"`
	Kind         string `json:"kind"`
	Category     string `json:"category"`
	Language     string `json:"language"`
	TimeLimitMin int    `json:"time_limit_min"`
	TutorMode    string `json:"tutor_mode"`
	RepoPath     string `json:"repo_path"`
	TestCommand  string `json:"test_command"`
	VideoURL     string `json:"video_url"`
	Difficulty   string `json:"difficulty"`
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

	kind := r.Kind
	if kind == "" {
		// Every exercise.json authored before kinds existed is a coding
		// exercise -- absence means coding, so none of them change.
		kind = KindCoding
	}
	if kind != KindCoding && kind != KindDesign {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid kind %q", path, r.Kind)
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
	languages, tutorModes := validLanguages, validTutorModes
	if kind == KindDesign {
		languages, tutorModes = validDesignLanguages, validDesignTutorModes
	}
	if !languages[r.Language] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid language %q for kind %q", path, r.Language, kind)
	}
	if r.TimeLimitMin <= 0 {
		return Exercise{}, fmt.Errorf("exercise: %s: time_limit_min must be > 0", path)
	}
	if !tutorModes[r.TutorMode] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid tutor_mode %q for kind %q", path, r.TutorMode, kind)
	}
	if r.RepoPath == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: repo_path is required", path)
	}
	if kind == KindDesign && r.TestCommand != "" {
		// A design session has nothing to run -- a test_command here is a
		// contradiction that would silently change submit semantics, so
		// fail loud at authoring time instead.
		return Exercise{}, fmt.Errorf("exercise: %s: test_command must be empty for a design exercise, got %q", path, r.TestCommand)
	}
	if kind == KindCoding && r.TestCommand == "" {
		return Exercise{}, fmt.Errorf("exercise: %s: test_command is required", path)
	}

	repoPath := r.RepoPath
	if !filepath.IsAbs(repoPath) {
		repoPath = filepath.Join(filepath.Dir(path), repoPath)
	}
	if info, err := os.Stat(repoPath); err != nil || !info.IsDir() {
		return Exercise{}, fmt.Errorf("exercise: %s: repo_path %q does not exist or is not a directory", path, repoPath)
	}

	if r.VideoURL != "" && !strings.HasPrefix(r.VideoURL, "https://") {
		return Exercise{}, fmt.Errorf("exercise: %s: video_url must be https, got %q", path, r.VideoURL)
	}

	if r.Difficulty != "" && !validDifficulties[r.Difficulty] {
		return Exercise{}, fmt.Errorf("exercise: %s: invalid difficulty %q (want easy|medium|hard)", path, r.Difficulty)
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
		Kind:         kind,
		Category:     r.Category,
		Language:     r.Language,
		TimeLimitMin: r.TimeLimitMin,
		TutorMode:    r.TutorMode,
		RepoPath:     repoPath,
		TestCommand:  r.TestCommand,
		VideoURL:     r.VideoURL,
		Difficulty:   r.Difficulty,
	}, nil
}
