package tracker

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func openTest(t *testing.T) *Tracker {
	t.Helper()
	dir := t.TempDir()
	tr, err := Open(filepath.Join(dir, "tracker.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { tr.Close() })
	return tr
}

func TestOpen_CreatesSchema(t *testing.T) {
	tr := openTest(t)

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts on fresh db: %v", err)
	}
	if len(attempts) != 0 {
		t.Fatalf("expected 0 attempts on fresh db, got %d", len(attempts))
	}
}

func TestOpen_MigratesPatternCategoryToDSA(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tracker.db")

	tr := func() *Tracker {
		tr, err := Open(path)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		return tr
	}()
	if _, err := tr.LogAttempt(Attempt{
		ExerciseID:   "two-pointers-01",
		Category:     "pattern",
		Language:     "go",
		Date:         "2026-07-08",
		TimeSpentMin: 10,
		Result:       ResultPass,
	}); err != nil {
		t.Fatalf("LogAttempt: %v", err)
	}
	tr.Close()

	tr2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer tr2.Close()

	attempts, err := tr2.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if attempts[0].Category != "dsa" {
		t.Errorf("Category = %q, want migrated %q", attempts[0].Category, "dsa")
	}
}

func TestLogAttempt_InsertsRow(t *testing.T) {
	tr := openTest(t)

	a := Attempt{
		ExerciseID:   "two-pointers-01",
		Category:     "pattern",
		Language:     "go",
		Date:         "2026-07-08",
		TimeSpentMin: 22.5,
		Result:       ResultPass,
		Notes:        "clean first try",
	}

	id, err := tr.LogAttempt(a)
	if err != nil {
		t.Fatalf("LogAttempt: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero id")
	}

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}

	got := attempts[0]
	if got.ID != id {
		t.Errorf("ID = %d, want %d", got.ID, id)
	}
	if got.ExerciseID != a.ExerciseID {
		t.Errorf("ExerciseID = %q, want %q", got.ExerciseID, a.ExerciseID)
	}
	if got.Category != a.Category {
		t.Errorf("Category = %q, want %q", got.Category, a.Category)
	}
	if got.Language != a.Language {
		t.Errorf("Language = %q, want %q", got.Language, a.Language)
	}
	if got.Date != a.Date {
		t.Errorf("Date = %q, want %q", got.Date, a.Date)
	}
	if got.TimeSpentMin != a.TimeSpentMin {
		t.Errorf("TimeSpentMin = %v, want %v", got.TimeSpentMin, a.TimeSpentMin)
	}
	if got.Result != a.Result {
		t.Errorf("Result = %q, want %q", got.Result, a.Result)
	}
	if got.Notes != a.Notes {
		t.Errorf("Notes = %q, want %q", got.Notes, a.Notes)
	}
}

func TestLogAttempt_RejectsInvalidResult(t *testing.T) {
	tr := openTest(t)

	a := Attempt{
		ExerciseID:   "x",
		Category:     "pattern",
		Language:     "go",
		Date:         "2026-07-08",
		TimeSpentMin: 1,
		Result:       "maybe",
	}

	if _, err := tr.LogAttempt(a); err == nil {
		t.Fatalf("expected error for invalid result %q, got nil", a.Result)
	}
}

func TestListAttempts_OrdersByIDAscending(t *testing.T) {
	tr := openTest(t)

	for _, id := range []string{"first", "second", "third"} {
		a := Attempt{
			ExerciseID:   id,
			Category:     "pattern",
			Language:     "go",
			Date:         "2026-07-08",
			TimeSpentMin: 1,
			Result:       ResultFail,
		}
		if _, err := tr.LogAttempt(a); err != nil {
			t.Fatalf("LogAttempt(%s): %v", id, err)
		}
	}

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(attempts))
	}
	want := []string{"first", "second", "third"}
	for i, w := range want {
		if attempts[i].ExerciseID != w {
			t.Errorf("attempts[%d].ExerciseID = %q, want %q", i, attempts[i].ExerciseID, w)
		}
	}
}

func TestOpen_MigratesGradeSummaryColumnOntoExistingDB(t *testing.T) {
	// A database created by an older binary has no grade_summary column;
	// Open must add it rather than failing or wiping rows.
	path := filepath.Join(t.TempDir(), "tracker.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		exercise_id TEXT NOT NULL, category TEXT NOT NULL,
		language TEXT NOT NULL, date TEXT NOT NULL,
		time_spent_min REAL NOT NULL,
		result TEXT NOT NULL CHECK (result IN ('pass', 'fail')),
		notes TEXT NOT NULL DEFAULT ''
	)`); err != nil {
		t.Fatalf("create legacy table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result, notes)
		VALUES ('two-sum-01', 'dsa', 'go', '2026-07-01', 12, 'pass', 'old row')`); err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}
	db.Close()

	tr, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer tr.Close()
	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 || attempts[0].Notes != "old row" || attempts[0].GradeSummary != "" {
		t.Errorf("legacy row after migration = %+v, want preserved with empty GradeSummary", attempts)
	}
}

func TestLogAttempt_RoundTripsGradeSummary(t *testing.T) {
	tr, err := Open(filepath.Join(t.TempDir(), "tracker.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer tr.Close()

	want := "VERDICT: fail\n1. Estimates: missing."
	if _, err := tr.LogAttempt(Attempt{
		ExerciseID: "url-shortener-01-interviewer", Category: "system-design",
		Language: "interviewer", Date: "2026-07-15", TimeSpentMin: 40,
		Result: ResultFail, GradeSummary: want,
	}); err != nil {
		t.Fatalf("LogAttempt: %v", err)
	}
	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 || attempts[0].GradeSummary != want {
		t.Errorf("GradeSummary = %q, want %q", attempts[0].GradeSummary, want)
	}
}
