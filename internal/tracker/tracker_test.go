package tracker

import (
	"database/sql"
	"errors"
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

func TestOpen_FreshDBLandsAtHighestVersionWithFullSchema(t *testing.T) {
	tr := openTest(t)

	v, err := currentVersion(tr.db)
	if err != nil {
		t.Fatalf("currentVersion: %v", err)
	}
	want := migrations[len(migrations)-1].version
	if v != want {
		t.Errorf("schema_version = %d, want latest %d", v, want)
	}

	for _, col := range []string{"grade_summary", "hints_used", "tutor_mode", "model", "turns"} {
		has, err := hasColumn(tr.db, "attempts", col)
		if err != nil {
			t.Fatalf("hasColumn(%s): %v", col, err)
		}
		if !has {
			t.Errorf("fresh schema missing column %q", col)
		}
	}

	// The CHECK constraint must already allow gave-up on a fresh db --
	// migration 002 is not something a later Open should still owe us.
	if _, err := tr.LogAttempt(Attempt{
		ExerciseID: "x", Category: "dsa", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 1, Result: ResultGaveUp,
	}); err != nil {
		t.Errorf("LogAttempt with gave-up on fresh db: %v", err)
	}
}

// TestMigrate_OldSchemaConvergesWithZeroDataLoss seeds a database with the
// literal pre-migration-runner schema (raw CREATE TABLE, no schema_version,
// no grade_summary, CHECK allows only pass/fail) and rows including a
// category="pattern" row, exactly like a real database that has existed
// since before this migration runner shipped. Opening it must reclassify
// pattern -> dsa, add grade_summary, and lose nothing.
func TestMigrate_OldSchemaConvergesWithZeroDataLoss(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tracker.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE attempts (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		exercise_id     TEXT NOT NULL,
		category        TEXT NOT NULL,
		language        TEXT NOT NULL,
		date            TEXT NOT NULL,
		time_spent_min  REAL NOT NULL,
		result          TEXT NOT NULL CHECK (result IN ('pass', 'fail')),
		notes           TEXT NOT NULL DEFAULT ''
	)`); err != nil {
		t.Fatalf("create legacy table: %v", err)
	}

	seedRows := []struct {
		exerciseID, category, language, date, result, notes string
		timeSpent                                           float64
	}{
		{exerciseID: "two-pointers-01", category: "pattern", language: "go", date: "2026-06-01", timeSpent: 10, result: "pass", notes: "first try, clean"},
		{exerciseID: "two-sum-01", category: "dsa", language: "python", date: "2026-06-02", timeSpent: 20, result: "fail", notes: "off by one"},
		{exerciseID: "lru-cache-01", category: "pattern", language: "go", date: "2026-06-03", timeSpent: 30, result: "fail", notes: ""},
	}
	for _, r := range seedRows {
		if _, err := db.Exec(
			`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result, notes)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			r.exerciseID, r.category, r.language, r.date, r.timeSpent, r.result, r.notes,
		); err != nil {
			t.Fatalf("seed insert %s: %v", r.exerciseID, err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close raw: %v", err)
	}

	tr, err := Open(path)
	if err != nil {
		t.Fatalf("Open on legacy db: %v", err)
	}
	defer tr.Close()

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != len(seedRows) {
		t.Fatalf("zero data loss violated: got %d attempts, want %d", len(attempts), len(seedRows))
	}

	for i, want := range seedRows {
		got := attempts[i]
		wantCategory := want.category
		if wantCategory == "pattern" {
			wantCategory = "dsa" // the regression under test
		}
		if got.ExerciseID != want.exerciseID {
			t.Errorf("row %d ExerciseID = %q, want %q", i, got.ExerciseID, want.exerciseID)
		}
		if got.Category != wantCategory {
			t.Errorf("row %d Category = %q, want %q", i, got.Category, wantCategory)
		}
		if got.Language != want.language {
			t.Errorf("row %d Language = %q, want %q", i, got.Language, want.language)
		}
		if got.Date != want.date {
			t.Errorf("row %d Date = %q, want %q", i, got.Date, want.date)
		}
		if got.TimeSpentMin != want.timeSpent {
			t.Errorf("row %d TimeSpentMin = %v, want %v", i, got.TimeSpentMin, want.timeSpent)
		}
		if got.Result != want.result {
			t.Errorf("row %d Result = %q, want %q", i, got.Result, want.result)
		}
		if got.Notes != want.notes {
			t.Errorf("row %d Notes = %q, want %q", i, got.Notes, want.notes)
		}
		if got.GradeSummary != "" {
			t.Errorf("row %d GradeSummary = %q, want empty default", i, got.GradeSummary)
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

func TestMigrate_RunningTwiceAppliesNothingTheSecondTime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tracker.db")

	tr, err := Open(path)
	if err != nil {
		t.Fatalf("Open (first): %v", err)
	}
	if _, err := tr.LogAttempt(Attempt{
		ExerciseID: "probe", Category: "pattern", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 1, Result: ResultPass,
	}); err != nil {
		t.Fatalf("LogAttempt: %v", err)
	}
	before, err := currentVersion(tr.db)
	if err != nil {
		t.Fatalf("currentVersion before: %v", err)
	}
	tr.Close()

	tr2, err := Open(path)
	if err != nil {
		t.Fatalf("Open (second): %v", err)
	}
	defer tr2.Close()

	after, err := currentVersion(tr2.db)
	if err != nil {
		t.Fatalf("currentVersion after: %v", err)
	}
	if after != before {
		t.Errorf("version changed on second Open: before=%d after=%d", before, after)
	}

	attempts, err := tr2.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	// The pattern->dsa backfill is a one-time historical migration, not
	// steady-state behavior -- a row logged as "pattern" after the
	// database is already on the latest schema must NOT get silently
	// rewritten by a later Open. Reopening must be a true no-op.
	if attempts[0].Category != "pattern" {
		t.Errorf("second Open mutated data (Category = %q); migrations must not re-run", attempts[0].Category)
	}
}

// TestMigrate_FailingMigrationRollsBackAndLeavesVersionUnchanged swaps in a
// migration that mutates the schema and then fails, and checks that
// neither the mutation nor the version bump survive.
func TestMigrate_FailingMigrationRollsBackAndLeavesVersionUnchanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tracker.db")

	origMigrations := migrations
	t.Cleanup(func() { migrations = origMigrations })

	// Seed a real database at version 1 only.
	migrations = []migration{{1, migrate001}}
	tr, err := Open(path)
	if err != nil {
		t.Fatalf("Open (seed v1): %v", err)
	}
	tr.Close()

	// Now attempt a v2 that mutates the schema before failing.
	migrations = []migration{
		{1, migrate001},
		{2, func(tx *sql.Tx) error {
			if _, err := tx.Exec(`ALTER TABLE attempts ADD COLUMN should_not_persist TEXT`); err != nil {
				return err
			}
			return errors.New("boom: simulated migration failure")
		}},
	}

	if _, err := Open(path); err == nil {
		t.Fatalf("expected error from failing migration, got nil")
	}

	// Inspect the persisted state directly, bypassing our Open (which
	// would just try to migrate again).
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("raw open: %v", err)
	}
	defer db.Close()

	v, err := currentVersion(db)
	if err != nil {
		t.Fatalf("currentVersion: %v", err)
	}
	if v != 1 {
		t.Errorf("version after failed migration = %d, want unchanged at 1", v)
	}

	has, err := hasColumn(db, "attempts", "should_not_persist")
	if err != nil {
		t.Fatalf("hasColumn: %v", err)
	}
	if has {
		t.Errorf("should_not_persist column leaked out of the rolled-back transaction")
	}
}

func TestMigrate002_RebuildsCheckConstraintAtTheSQLLevel(t *testing.T) {
	// Bypasses LogAttempt's own Go-side guard so this proves the CHECK
	// constraint on the rebuilt table itself, not just the wrapper.
	tr := openTest(t)

	if _, err := tr.db.Exec(
		`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result)
		 VALUES ('x', 'dsa', 'go', '2026-07-17', 1, 'gave-up')`,
	); err != nil {
		t.Errorf("SQL-level insert of gave-up rejected: %v", err)
	}

	if _, err := tr.db.Exec(
		`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result)
		 VALUES ('x', 'dsa', 'go', '2026-07-17', 1, 'bogus')`,
	); err == nil {
		t.Errorf("SQL-level insert of an unknown result value should have violated the CHECK constraint")
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

func TestLogAttempt_AcceptsGaveUpResult(t *testing.T) {
	tr := openTest(t)

	id, err := tr.LogAttempt(Attempt{
		ExerciseID: "two-pointers-01", Category: "dsa", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 30, Result: ResultGaveUp,
	})
	if err != nil {
		t.Fatalf("LogAttempt(gave-up): %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero id")
	}

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 || attempts[0].Result != ResultGaveUp {
		t.Errorf("attempts = %+v, want a single gave-up result", attempts)
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

func TestLogAttempt_RoundTripsNewNullableColumns(t *testing.T) {
	tr := openTest(t)

	hints := 2
	mode := "socratic"
	model := "claude-sonnet-5"
	turns := 7
	full := Attempt{
		ExerciseID: "two-sum-01", Category: "dsa", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 15, Result: ResultPass,
		HintsUsed: &hints, TutorMode: &mode, Model: &model, Turns: &turns,
	}
	empty := Attempt{
		ExerciseID: "two-sum-02", Category: "dsa", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 5, Result: ResultFail,
	}

	if _, err := tr.LogAttempt(full); err != nil {
		t.Fatalf("LogAttempt(full): %v", err)
	}
	if _, err := tr.LogAttempt(empty); err != nil {
		t.Fatalf("LogAttempt(empty): %v", err)
	}

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 attempts, got %d", len(attempts))
	}

	got := attempts[0]
	if got.HintsUsed == nil || *got.HintsUsed != hints {
		t.Errorf("HintsUsed = %v, want %d", got.HintsUsed, hints)
	}
	if got.TutorMode == nil || *got.TutorMode != mode {
		t.Errorf("TutorMode = %v, want %q", got.TutorMode, mode)
	}
	if got.Model == nil || *got.Model != model {
		t.Errorf("Model = %v, want %q", got.Model, model)
	}
	if got.Turns == nil || *got.Turns != turns {
		t.Errorf("Turns = %v, want %d", got.Turns, turns)
	}

	got2 := attempts[1]
	if got2.HintsUsed != nil {
		t.Errorf("HintsUsed = %v, want nil", got2.HintsUsed)
	}
	if got2.TutorMode != nil {
		t.Errorf("TutorMode = %v, want nil", got2.TutorMode)
	}
	if got2.Model != nil {
		t.Errorf("Model = %v, want nil", got2.Model)
	}
	if got2.Turns != nil {
		t.Errorf("Turns = %v, want nil", got2.Turns)
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

func TestListAttemptsFor_FiltersAndOrdersNewestFirst(t *testing.T) {
	tr := openTest(t)

	log := func(exerciseID, date string) {
		t.Helper()
		if _, err := tr.LogAttempt(Attempt{
			ExerciseID: exerciseID, Category: "dsa", Language: "go",
			Date: date, TimeSpentMin: 1, Result: ResultPass,
		}); err != nil {
			t.Fatalf("LogAttempt(%s, %s): %v", exerciseID, date, err)
		}
	}

	log("two-sum-01", "2026-07-01")
	log("lru-cache-01", "2026-07-05") // different exercise, interleaved
	log("two-sum-01", "2026-07-10")
	log("two-sum-01", "2026-07-15")

	got, err := tr.ListAttemptsFor("two-sum-01")
	if err != nil {
		t.Fatalf("ListAttemptsFor: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 attempts for two-sum-01, got %d", len(got))
	}
	wantDates := []string{"2026-07-15", "2026-07-10", "2026-07-01"} // newest first
	for i, want := range wantDates {
		if got[i].Date != want {
			t.Errorf("attempt[%d].Date = %q, want %q (newest-first order)", i, got[i].Date, want)
		}
		if got[i].ExerciseID != "two-sum-01" {
			t.Errorf("attempt[%d].ExerciseID = %q, want two-sum-01 (filter leaked)", i, got[i].ExerciseID)
		}
	}
}

func TestListAttemptsFor_UnknownExerciseReturnsEmpty(t *testing.T) {
	tr := openTest(t)

	if _, err := tr.LogAttempt(Attempt{
		ExerciseID: "two-sum-01", Category: "dsa", Language: "go",
		Date: "2026-07-17", TimeSpentMin: 1, Result: ResultPass,
	}); err != nil {
		t.Fatalf("LogAttempt: %v", err)
	}

	got, err := tr.ListAttemptsFor("does-not-exist")
	if err != nil {
		t.Fatalf("ListAttemptsFor: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 attempts, got %d", len(got))
	}
}
