// Package tracker reads and writes the SQLite attempts log at data/tracker.db.
// Single flat table of attempts — per-category tables would buy nothing
// a WHERE clause doesn't, and scheduling lives in internal/catalog
// (due.go's review windows) rather than in the schema.
package tracker

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	ResultPass = "pass"
	ResultFail = "fail"
)

const schema = `
CREATE TABLE IF NOT EXISTS attempts (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	exercise_id     TEXT NOT NULL,
	category        TEXT NOT NULL,
	language        TEXT NOT NULL,
	date            TEXT NOT NULL,
	time_spent_min  REAL NOT NULL,
	result          TEXT NOT NULL CHECK (result IN ('pass', 'fail')),
	notes           TEXT NOT NULL DEFAULT '',
	grade_summary   TEXT NOT NULL DEFAULT ''
);
`

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
}

// Tracker is a handle to the attempts SQLite database.
type Tracker struct {
	db *sql.DB
}

// Open opens (creating if needed) the SQLite database at path and ensures
// the attempts table exists.
func Open(path string) (*Tracker, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("tracker: open %s: %w", path, err)
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("tracker: migrate schema: %w", err)
	}
	// The "pattern" category was renamed to "dsa" — reclassify any
	// already-logged attempts so Stats keeps grouping them together
	// instead of splitting into two buckets under the old and new names.
	if _, err := db.Exec(`UPDATE attempts SET category = 'dsa' WHERE category = 'pattern'`); err != nil {
		db.Close()
		return nil, fmt.Errorf("tracker: migrate pattern category: %w", err)
	}
	// grade_summary arrived after databases already existed in the wild;
	// CREATE TABLE IF NOT EXISTS won't touch those, so probe and ALTER.
	if _, err := db.Exec(`SELECT grade_summary FROM attempts LIMIT 0`); err != nil {
		if _, err := db.Exec(`ALTER TABLE attempts ADD COLUMN grade_summary TEXT NOT NULL DEFAULT ''`); err != nil {
			db.Close()
			return nil, fmt.Errorf("tracker: migrate grade_summary column: %w", err)
		}
	}
	return &Tracker{db: db}, nil
}

// Close releases the underlying database handle.
func (t *Tracker) Close() error {
	return t.db.Close()
}

// LogAttempt inserts a new attempt row and returns its id.
func (t *Tracker) LogAttempt(a Attempt) (int64, error) {
	if a.Result != ResultPass && a.Result != ResultFail {
		return 0, fmt.Errorf("tracker: invalid result %q (want %q or %q)", a.Result, ResultPass, ResultFail)
	}

	res, err := t.db.Exec(
		`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result, notes, grade_summary)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ExerciseID, a.Category, a.Language, a.Date, a.TimeSpentMin, a.Result, a.Notes, a.GradeSummary,
	)
	if err != nil {
		return 0, fmt.Errorf("tracker: insert attempt: %w", err)
	}
	return res.LastInsertId()
}

// ListAttempts returns all attempts ordered by id ascending.
func (t *Tracker) ListAttempts() ([]Attempt, error) {
	rows, err := t.db.Query(
		`SELECT id, exercise_id, category, language, date, time_spent_min, result, notes, grade_summary
		 FROM attempts ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("tracker: list attempts: %w", err)
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		var a Attempt
		if err := rows.Scan(&a.ID, &a.ExerciseID, &a.Category, &a.Language, &a.Date, &a.TimeSpentMin, &a.Result, &a.Notes, &a.GradeSummary); err != nil {
			return nil, fmt.Errorf("tracker: scan attempt: %w", err)
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}
