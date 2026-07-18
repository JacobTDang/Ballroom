// Package tracker reads and writes the SQLite attempts log at data/tracker.db.
// Single flat table of attempts — per-category tables would buy nothing
// a WHERE clause doesn't, and scheduling lives in internal/catalog
// (due.go's review windows) rather than in the schema.
//
// Schema changes go through the ordered migrations slice below: each
// migration runs once, inside its own transaction, recording its version
// in schema_version only on success. This replaced an ad-hoc
// CREATE-IF-NOT-EXISTS-plus-probe-and-ALTER approach that couldn't
// express a CHECK constraint rebuild -- SQLite has no ALTER COLUMN or
// DROP CONSTRAINT, so widening a CHECK requires rebuilding the table.
package tracker

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	ResultPass   = "pass"
	ResultFail   = "fail"
	ResultGaveUp = "gave-up"
)

// Attempt is one row of the attempts table.
type Attempt struct {
	ID           int64
	ExerciseID   string
	Category     string
	Language     string
	Date         string
	TimeSpentMin float64
	Result       string
	Notes        string
	// GradeSummary is the design grader's per-dimension assessment (see
	// tutor.GradeDesign) -- empty for coding attempts and self-assessed
	// design attempts. Persisted so Stats can aggregate weak dimensions
	// after the session's workspace is gone.
	GradeSummary string
	// HintsUsed, TutorMode, Model, and Turns describe tutor-assisted
	// sessions and are nil for a plain submit the tutor never touched.
	// Nullable in the schema (migration 002), so pointers here rather
	// than zero-valuing an int or string that could otherwise mean
	// "zero hints used" as easily as "no tutor session at all".
	HintsUsed *int
	TutorMode *string
	Model     *string
	Turns     *int
}

// Tracker is a handle to the attempts SQLite database.
type Tracker struct {
	db *sql.DB
}

// migration is one step in the schema's history. up runs inside a single
// transaction; returning an error rolls back everything it did,
// including the schema_version bump, so a failed migration never leaves
// the database at a version whose changes didn't fully land.
type migration struct {
	version int
	up      func(tx *sql.Tx) error
}

// migrations is the ordered schema history. Append new steps here --
// never edit a step that has already shipped, since real databases may
// already be past it and its job is to describe what already ran.
var migrations = []migration{
	{1, migrate001},
	{2, migrate002},
}

// migrate001 is the schema as it already existed, in every database in
// the wild, the moment this migration runner shipped: the attempts
// table plus the two ad-hoc steps that used to run unconditionally on
// every Open (the pattern -> dsa category rename and the grade_summary
// column) folded in as one-time steps. A fresh database and a database
// upgraded from the pre-migration-runner shape both run this exact
// sequence, so both converge on the same end state.
func migrate001(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS attempts (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			exercise_id     TEXT NOT NULL,
			category        TEXT NOT NULL,
			language        TEXT NOT NULL,
			date            TEXT NOT NULL,
			time_spent_min  REAL NOT NULL,
			result          TEXT NOT NULL CHECK (result IN ('pass', 'fail')),
			notes           TEXT NOT NULL DEFAULT ''
		)`); err != nil {
		return fmt.Errorf("create attempts: %w", err)
	}

	// grade_summary arrived after databases already existed in the wild;
	// CREATE TABLE IF NOT EXISTS won't touch those, so probe and ALTER.
	has, err := hasColumn(tx, "attempts", "grade_summary")
	if err != nil {
		return err
	}
	if !has {
		if _, err := tx.Exec(`ALTER TABLE attempts ADD COLUMN grade_summary TEXT NOT NULL DEFAULT ''`); err != nil {
			return fmt.Errorf("add grade_summary: %w", err)
		}
	}

	// The "pattern" category was renamed to "dsa" -- reclassify any
	// already-logged attempts so Stats keeps grouping them together
	// instead of splitting into two buckets under the old and new names.
	// A one-time historical cleanup: it runs here, as part of reaching
	// version 1, and never again.
	if _, err := tx.Exec(`UPDATE attempts SET category = 'dsa' WHERE category = 'pattern'`); err != nil {
		return fmt.Errorf("backfill pattern category: %w", err)
	}
	return nil
}

// migrate002 adds tutor-session metadata columns and widens the result
// CHECK constraint to allow "gave-up". SQLite cannot ALTER a CHECK
// constraint in place, so this uses the standard rebuild dance: build
// the new shape under a temporary name, copy every row across, drop the
// old table, then rename the new one into place.
func migrate002(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		CREATE TABLE attempts_new (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			exercise_id     TEXT NOT NULL,
			category        TEXT NOT NULL,
			language        TEXT NOT NULL,
			date            TEXT NOT NULL,
			time_spent_min  REAL NOT NULL,
			result          TEXT NOT NULL CHECK (result IN ('pass', 'fail', 'gave-up')),
			notes           TEXT NOT NULL DEFAULT '',
			grade_summary   TEXT NOT NULL DEFAULT '',
			hints_used      INTEGER,
			tutor_mode      TEXT,
			model           TEXT,
			turns           INTEGER
		)`); err != nil {
		return fmt.Errorf("create attempts_new: %w", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO attempts_new (id, exercise_id, category, language, date, time_spent_min, result, notes, grade_summary)
		SELECT id, exercise_id, category, language, date, time_spent_min, result, notes, grade_summary FROM attempts
	`); err != nil {
		return fmt.Errorf("copy attempts into attempts_new: %w", err)
	}
	if _, err := tx.Exec(`DROP TABLE attempts`); err != nil {
		return fmt.Errorf("drop old attempts: %w", err)
	}
	if _, err := tx.Exec(`ALTER TABLE attempts_new RENAME TO attempts`); err != nil {
		return fmt.Errorf("rename attempts_new to attempts: %w", err)
	}
	return nil
}

// queryer is satisfied by both *sql.DB and *sql.Tx, so hasColumn works
// during a migration (which only has a *sql.Tx) and from tests that want
// to assert on the final schema (which only have a *sql.DB).
type queryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// hasColumn reports whether table has a column named col. SQLite has no
// information_schema; PRAGMA table_info is the documented way to
// introspect a table's columns. table is always an internal literal, so
// building the PRAGMA with Sprintf carries no injection risk -- SQLite
// also doesn't accept bound parameters inside a PRAGMA.
func hasColumn(q queryer, table, col string) (bool, error) {
	rows, err := q.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, table))
	if err != nil {
		return false, fmt.Errorf("table_info(%s): %w", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false, fmt.Errorf("scan table_info(%s): %w", table, err)
		}
		if name == col {
			return true, nil
		}
	}
	return false, rows.Err()
}

// currentVersion reads the single row in schema_version, or 0 if the
// table is empty (a database that has never run a migration).
func currentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow(`SELECT version FROM schema_version LIMIT 1`).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read schema_version: %w", err)
	}
	return version, nil
}

// applyMigration runs one migration's up function and records its
// version, all inside one transaction -- either both happen or neither
// does.
func applyMigration(db *sql.DB, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := m.up(tx); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`DELETE FROM schema_version`); err != nil {
		tx.Rollback()
		return fmt.Errorf("clear schema_version: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES (?)`, m.version); err != nil {
		tx.Rollback()
		return fmt.Errorf("record schema_version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// migrate brings db up to the latest schema version, running each
// pending migration in its own transaction in order. Idempotent: a
// database already at the latest version runs nothing.
func migrate(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`); err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}

	current, err := currentVersion(db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version <= current {
			continue
		}
		if err := applyMigration(db, m); err != nil {
			return fmt.Errorf("migration %d: %w", m.version, err)
		}
		current = m.version
	}
	return nil
}

// Open opens (creating if needed) the SQLite database at path and brings
// its schema up to the latest version, running any pending migrations.
func Open(path string) (*Tracker, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("tracker: open %s: %w", path, err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("tracker: %w", err)
	}
	return &Tracker{db: db}, nil
}

// Close releases the underlying database handle.
func (t *Tracker) Close() error {
	return t.db.Close()
}

// LogAttempt inserts a new attempt row and returns its id.
func (t *Tracker) LogAttempt(a Attempt) (int64, error) {
	switch a.Result {
	case ResultPass, ResultFail, ResultGaveUp:
	default:
		return 0, fmt.Errorf("tracker: invalid result %q (want %q, %q, or %q)", a.Result, ResultPass, ResultFail, ResultGaveUp)
	}

	res, err := t.db.Exec(
		`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result, notes, grade_summary, hints_used, tutor_mode, model, turns)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ExerciseID, a.Category, a.Language, a.Date, a.TimeSpentMin, a.Result, a.Notes, a.GradeSummary,
		a.HintsUsed, a.TutorMode, a.Model, a.Turns,
	)
	if err != nil {
		return 0, fmt.Errorf("tracker: insert attempt: %w", err)
	}
	return res.LastInsertId()
}

// attemptColumns is shared by ListAttempts and ListAttemptsFor so the
// SELECT list and the Scan order below can't drift apart.
const attemptColumns = `id, exercise_id, category, language, date, time_spent_min, result, notes, grade_summary, hints_used, tutor_mode, model, turns`

func scanAttempt(rows *sql.Rows) (Attempt, error) {
	var a Attempt
	err := rows.Scan(
		&a.ID, &a.ExerciseID, &a.Category, &a.Language, &a.Date, &a.TimeSpentMin, &a.Result, &a.Notes, &a.GradeSummary,
		&a.HintsUsed, &a.TutorMode, &a.Model, &a.Turns,
	)
	return a, err
}

// ListAttempts returns all attempts ordered by id ascending.
func (t *Tracker) ListAttempts() ([]Attempt, error) {
	rows, err := t.db.Query(`SELECT ` + attemptColumns + ` FROM attempts ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("tracker: list attempts: %w", err)
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		a, err := scanAttempt(rows)
		if err != nil {
			return nil, fmt.Errorf("tracker: scan attempt: %w", err)
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}

// ListAttemptsFor returns every attempt logged against exerciseID,
// newest first (highest id first -- attempts are never edited or
// reordered after insertion, so insertion order is chronological order).
// Used by the Stats drill-down to show one exercise's attempt history.
func (t *Tracker) ListAttemptsFor(exerciseID string) ([]Attempt, error) {
	rows, err := t.db.Query(`SELECT `+attemptColumns+` FROM attempts WHERE exercise_id = ? ORDER BY id DESC`, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("tracker: list attempts for %s: %w", exerciseID, err)
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		a, err := scanAttempt(rows)
		if err != nil {
			return nil, fmt.Errorf("tracker: scan attempt: %w", err)
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}
