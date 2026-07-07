// Package tracker reads and writes the SQLite attempts log at data/tracker.db.
// Single flat table, matching interview_prep_mvp_spec.md Section 3.5 — no
// per-category tables, no spaced-repetition scheduling for MVP.
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
	notes           TEXT NOT NULL DEFAULT ''
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
		`INSERT INTO attempts (exercise_id, category, language, date, time_spent_min, result, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.ExerciseID, a.Category, a.Language, a.Date, a.TimeSpentMin, a.Result, a.Notes,
	)
	if err != nil {
		return 0, fmt.Errorf("tracker: insert attempt: %w", err)
	}
	return res.LastInsertId()
}

// ListAttempts returns all attempts ordered by id ascending.
func (t *Tracker) ListAttempts() ([]Attempt, error) {
	rows, err := t.db.Query(
		`SELECT id, exercise_id, category, language, date, time_spent_min, result, notes
		 FROM attempts ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("tracker: list attempts: %w", err)
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		var a Attempt
		if err := rows.Scan(&a.ID, &a.ExerciseID, &a.Category, &a.Language, &a.Date, &a.TimeSpentMin, &a.Result, &a.Notes); err != nil {
			return nil, fmt.Errorf("tracker: scan attempt: %w", err)
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}
