package exercise

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeExercise(t *testing.T, dir string, fields map[string]any) string {
	t.Helper()
	base := map[string]any{
		"id":             "two-pointers-01",
		"title":          "Two Sum II",
		"category":       "dsa",
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

// writeDesignExercise is writeExercise with a valid kind:"design" base --
// design exercises use the language slot for session style, a design
// tutor mode, and no test command (there are no tests to run).
func writeDesignExercise(t *testing.T, dir string, fields map[string]any) string {
	t.Helper()
	base := map[string]any{
		"id":             "url-shortener-01-coach",
		"problem_id":     "url-shortener-01",
		"title":          "Design Pastebin / Bit.ly",
		"kind":           "design",
		"category":       "system-design",
		"language":       "coach",
		"time_limit_min": 90,
		"tutor_mode":     "design-coach",
		"repo_path":      "./repo",
		"test_command":   "",
	}
	for k, v := range fields {
		base[k] = v
	}
	return writeExercise(t, dir, base)
}

func TestLoad_KindDefaultsToCodingWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	path := writeExercise(t, dir, nil)

	ex, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ex.Kind != KindCoding {
		t.Errorf("Kind = %q, want %q for an exercise.json with no kind field (all 460 existing files)", ex.Kind, KindCoding)
	}
}

func TestLoad_ValidDesignExercise(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	path := writeDesignExercise(t, dir, nil)

	ex, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ex.Kind != KindDesign {
		t.Errorf("Kind = %q, want %q", ex.Kind, KindDesign)
	}
	if ex.Category != CategorySystemDesign {
		t.Errorf("Category = %q, want %q", ex.Category, CategorySystemDesign)
	}
	if ex.Language != LanguageCoach {
		t.Errorf("Language = %q, want %q", ex.Language, LanguageCoach)
	}
	if ex.TutorMode != TutorModeDesignCoach {
		t.Errorf("TutorMode = %q, want %q", ex.TutorMode, TutorModeDesignCoach)
	}
	if ex.TestCommand != "" {
		t.Errorf("TestCommand = %q, want empty for a design exercise", ex.TestCommand)
	}
}

func TestLoad_DesignValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		fields  map[string]any
		wantErr string // substring; "" = must load cleanly
	}{
		{"interviewer variant loads", map[string]any{
			"id": "url-shortener-01-interviewer", "language": "interviewer",
			"tutor_mode": "interviewer", "time_limit_min": 45,
		}, ""},
		{"behavioral interviewer variant loads", map[string]any{
			"id": "disagreement-01-interviewer", "category": "behavioral", "language": "interviewer",
			"tutor_mode": "behavioral-interviewer", "time_limit_min": 30,
		}, ""},
		{"story coach variant loads", map[string]any{
			"id": "disagreement-01-coach", "category": "behavioral", "language": "coach",
			"tutor_mode": "story-coach", "time_limit_min": 45,
		}, ""},
		{"coding language rejected for design", map[string]any{"language": "python"},
			"invalid language"},
		{"coding tutor mode rejected for design", map[string]any{"tutor_mode": "hints-first"},
			"invalid tutor_mode"},
		{"test_command must be empty for design", map[string]any{"test_command": "go test ./..."},
			"test_command"},
		{"unknown kind rejected", map[string]any{"kind": "essay"},
			"invalid kind"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
				t.Fatalf("mkdir repo: %v", err)
			}
			path := writeDesignExercise(t, dir, tc.fields)
			_, err := Load(path)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("Load: %v, want success", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Load err = %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestLoad_DesignStylesRejectedForCodingKind(t *testing.T) {
	// The coach/interviewer pseudo-languages and design tutor modes are
	// kind-gated -- a coding exercise claiming them must fail loud.
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	path := writeExercise(t, dir, map[string]any{"language": "coach"})
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "invalid language") {
		t.Errorf("Load err = %v, want invalid-language error for a coding exercise with language coach", err)
	}

	path = writeExercise(t, dir, map[string]any{"tutor_mode": "interviewer"})
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "invalid tutor_mode") {
		t.Errorf("Load err = %v, want invalid-tutor_mode error for a coding exercise with interviewer mode", err)
	}

	path = writeExercise(t, dir, map[string]any{"tutor_mode": "behavioral-interviewer"})
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "invalid tutor_mode") {
		t.Errorf("Load err = %v, want invalid-tutor_mode error for a coding exercise with behavioral-interviewer mode", err)
	}
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
	if ex.Category != CategoryDSA {
		t.Errorf("Category = %q, want %q", ex.Category, CategoryDSA)
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

func TestLoad_AcceptsEveryNeetCodeCategory(t *testing.T) {
	categories := []string{
		CategoryArraysHashing,
		CategoryTwoPointers,
		CategorySlidingWindow,
		CategoryStack,
		CategoryBinarySearch,
		CategoryLinkedList,
		CategoryTrees,
		CategoryTries,
		CategoryHeap,
		CategoryBacktracking,
		CategoryGraphs,
		CategoryAdvancedGraphs,
		CategoryDP1D,
		CategoryDP2D,
		CategoryGreedy,
		CategoryIntervals,
		CategoryMathGeometry,
		CategoryBitManipulation,
	}
	for _, cat := range categories {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "repo"), 0o755); err != nil {
			t.Fatalf("mkdir repo: %v", err)
		}
		path := writeExercise(t, dir, map[string]any{"category": cat})

		ex, err := Load(path)
		if err != nil {
			t.Errorf("Load with category %q: %v", cat, err)
			continue
		}
		if ex.Category != cat {
			t.Errorf("Category = %q, want %q", ex.Category, cat)
		}
	}
}

func TestLoad_ProblemIDParsed(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	path := writeExercise(t, dir, map[string]any{"id": "two-pointers-01-go", "problem_id": "two-pointers-01"})

	ex, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ex.ProblemID != "two-pointers-01" {
		t.Errorf("ProblemID = %q, want %q", ex.ProblemID, "two-pointers-01")
	}
}

func TestLoad_ProblemIDDefaultsToIDWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "repo"), 0o755)
	// no problem_id field at all — older/standalone exercises without
	// language siblings shouldn't need to add one just to load.
	path := writeExercise(t, dir, map[string]any{"id": "standalone-01"})

	ex, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ex.ProblemID != "standalone-01" {
		t.Errorf("ProblemID = %q, want it to default to ID %q", ex.ProblemID, "standalone-01")
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
