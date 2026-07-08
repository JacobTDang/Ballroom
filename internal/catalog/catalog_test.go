package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsi removes color escape codes so content assertions can check
// substance without being coupled to exactly how something is styled.
func stripAnsi(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func writeExercise(t *testing.T, exercisesDir, id string, fields map[string]any) {
	t.Helper()
	dir := filepath.Join(exercisesDir, id)
	repoDir := filepath.Join(dir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	base := map[string]any{
		"id":             id,
		"title":          "Title for " + id,
		"category":       "pattern",
		"language":       "go",
		"time_limit_min": 20,
		"tutor_mode":     "hints-first",
		"repo_path":      "./repo",
		"test_command":   "true",
	}
	for k, v := range fields {
		base[k] = v
	}

	data, err := json.Marshal(base)
	if err != nil {
		t.Fatalf("marshal exercise.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "exercise.json"), data, 0o644); err != nil {
		t.Fatalf("write exercise.json: %v", err)
	}
}

func fakeExercise(id, category, language, title string) exercise.Exercise {
	return exercise.Exercise{
		ID:           id,
		Title:        title,
		Category:     category,
		Language:     language,
		TimeLimitMin: 20,
		TutorMode:    exercise.TutorModeHintsFirst,
		RepoPath:     "/fake/repo",
		TestCommand:  "true",
	}
}

func testConfig(t *testing.T) config.Config {
	t.Helper()
	dir := t.TempDir()
	exercisesDir := filepath.Join(dir, "exercises")
	if err := os.MkdirAll(exercisesDir, 0o755); err != nil {
		t.Fatalf("mkdir exercises: %v", err)
	}
	return config.Config{
		ExercisesDir: exercisesDir,
		DBPath:       filepath.Join(dir, "tracker.db"),
	}
}

func TestList_ReturnsExercisesSkippingTemplate(t *testing.T) {
	cfg := testConfig(t)
	writeExercise(t, cfg.ExercisesDir, "two-pointers-01", nil)
	writeExercise(t, cfg.ExercisesDir, "cpp-debug-01", map[string]any{"category": "debug", "language": "cpp"})
	// _template has no exercise.json at all — matches the real repo's template dir.
	if err := os.MkdirAll(filepath.Join(cfg.ExercisesDir, "_template", "repo"), 0o755); err != nil {
		t.Fatalf("mkdir _template: %v", err)
	}

	statuses, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 exercises, got %d: %+v", len(statuses), statuses)
	}
	for _, s := range statuses {
		if s.Exercise.ID == "_template" {
			t.Error("_template should not appear in List results")
		}
	}
}

func TestList_SortsByCategoryThenID(t *testing.T) {
	cfg := testConfig(t)
	writeExercise(t, cfg.ExercisesDir, "z-ai-assisted-01", map[string]any{"category": "ai-assisted"})
	writeExercise(t, cfg.ExercisesDir, "a-pattern-02", map[string]any{"category": "pattern"})
	writeExercise(t, cfg.ExercisesDir, "b-pattern-01", map[string]any{"category": "pattern"})
	writeExercise(t, cfg.ExercisesDir, "y-debug-01", map[string]any{"category": "debug", "language": "cpp"})

	statuses, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(statuses) != 4 {
		t.Fatalf("expected 4 exercises, got %d", len(statuses))
	}

	var gotOrder []string
	for _, s := range statuses {
		gotOrder = append(gotOrder, s.Exercise.ID)
	}
	want := []string{"a-pattern-02", "b-pattern-01", "y-debug-01", "z-ai-assisted-01"}
	if strings.Join(gotOrder, ",") != strings.Join(want, ",") {
		t.Errorf("order = %v, want %v (category order pattern < debug < ai-assisted, alphabetical within category)", gotOrder, want)
	}
}

func TestList_ComputesAttemptsAndLastResult(t *testing.T) {
	cfg := testConfig(t)
	writeExercise(t, cfg.ExercisesDir, "two-pointers-01", nil)

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("tracker.Open: %v", err)
	}
	if _, err := tr.LogAttempt(tracker.Attempt{
		ExerciseID: "two-pointers-01", Category: "pattern", Language: "go",
		Date: "2026-07-08", TimeSpentMin: 5, Result: tracker.ResultFail,
	}); err != nil {
		t.Fatalf("LogAttempt 1: %v", err)
	}
	if _, err := tr.LogAttempt(tracker.Attempt{
		ExerciseID: "two-pointers-01", Category: "pattern", Language: "go",
		Date: "2026-07-08", TimeSpentMin: 3, Result: tracker.ResultPass,
	}); err != nil {
		t.Fatalf("LogAttempt 2: %v", err)
	}
	tr.Close()

	statuses, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(statuses))
	}
	if statuses[0].Attempts != 2 {
		t.Errorf("Attempts = %d, want 2", statuses[0].Attempts)
	}
	if statuses[0].LastResult != tracker.ResultPass {
		t.Errorf("LastResult = %q, want %q (most recent attempt)", statuses[0].LastResult, tracker.ResultPass)
	}
}

func TestList_NeverAttemptedShowsEmptyResult(t *testing.T) {
	cfg := testConfig(t)
	writeExercise(t, cfg.ExercisesDir, "two-pointers-01", nil)

	statuses, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(statuses))
	}
	if statuses[0].Attempts != 0 {
		t.Errorf("Attempts = %d, want 0", statuses[0].Attempts)
	}
	if statuses[0].LastResult != "" {
		t.Errorf("LastResult = %q, want empty (never attempted)", statuses[0].LastResult)
	}
}

func TestList_SkipsInvalidExerciseDirectories(t *testing.T) {
	cfg := testConfig(t)
	writeExercise(t, cfg.ExercisesDir, "two-pointers-01", nil)

	// A directory with no exercise.json at all should be skipped, not
	// cause the whole List() call to fail.
	if err := os.MkdirAll(filepath.Join(cfg.ExercisesDir, "broken-exercise"), 0o755); err != nil {
		t.Fatalf("mkdir broken-exercise: %v", err)
	}

	statuses, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 valid exercise (broken one skipped), got %d: %+v", len(statuses), statuses)
	}
}

func TestFormatTable_IncludesExerciseInfo(t *testing.T) {
	statuses := []ExerciseStatus{
		{Exercise: fakeExercise("two-pointers-01", "pattern", "go", "Two Sum II"), Attempts: 2, LastResult: tracker.ResultPass},
		{Exercise: fakeExercise("cpp-debug-01", "debug", "cpp", "Off-by-one"), Attempts: 0, LastResult: ""},
	}

	out := stripAnsi(FormatTable(statuses))

	for _, want := range []string{"two-pointers-01", "Two Sum II", "pattern", "go", "pass", "2 attempt"} {
		if !strings.Contains(out, want) {
			t.Errorf("table output missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(strings.ToLower(out), "not attempted") {
		t.Errorf("table output missing 'not attempted' for zero-attempt exercise:\n%s", out)
	}
}

func TestFormatSummary_ShowsPerCategoryCounts(t *testing.T) {
	statuses := []ExerciseStatus{
		{Exercise: fakeExercise("a", "pattern", "go", "A"), LastResult: tracker.ResultPass},
		{Exercise: fakeExercise("b", "debug", "cpp", "B"), LastResult: tracker.ResultFail},
		{Exercise: fakeExercise("c", "debug", "cpp", "C"), LastResult: ""},
	}

	out := stripAnsi(FormatSummary(statuses))

	if !strings.Contains(out, "pattern: 1/1") {
		t.Errorf("summary missing pattern 1/1:\n%s", out)
	}
	if !strings.Contains(out, "debug: 0/2") {
		t.Errorf("summary missing debug 0/2:\n%s", out)
	}
}
