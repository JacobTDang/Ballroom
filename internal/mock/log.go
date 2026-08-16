package mock

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Per-slot outcomes recorded in the sitting log.
const (
	OutcomeSolved    = "solved"    // a submit passed
	OutcomeAttempted = "attempted" // attempts logged, none passed
	OutcomeSkipped   = "skipped"   // never launched, or launched with no attempt
)

// Sitting is one mock run — a row of mocks.jsonl. Facts only: no
// synthetic CodeSignal-style score, which couldn't be calibrated
// honestly from outside.
type Sitting struct {
	StartedAt      string     `json:"started_at"`
	ProblemIDs     [4]string  `json:"problem_ids"`
	Outcomes       [4]string  `json:"outcomes"`
	MinutesPerSlot [4]float64 `json:"minutes_per_slot"`
	MinutesTotal   float64    `json:"minutes_total"`
	Completed      bool       `json:"completed"`
}

// Solved counts the sitting's solved slots.
func (s Sitting) Solved() int {
	n := 0
	for _, o := range s.Outcomes {
		if o == OutcomeSolved {
			n++
		}
	}
	return n
}

func logPath(dataDir string) string { return filepath.Join(dataDir, "mocks.jsonl") }

// AppendSitting appends one sitting as a JSON line.
func AppendSitting(dataDir string, s Sitting) error {
	f, err := os.OpenFile(logPath(dataDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("mock: open sitting log: %w", err)
	}
	defer f.Close()
	line, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("mock: encode sitting: %w", err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("mock: append sitting: %w", err)
	}
	return nil
}

// ListSittings reads the whole log, oldest first. A missing file is an
// empty history; a malformed line is an error — a half-readable log
// silently dropping rows would misreport history.
func ListSittings(dataDir string) ([]Sitting, error) {
	f, err := os.Open(logPath(dataDir))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("mock: open sitting log: %w", err)
	}
	defer f.Close()

	var out []Sitting
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var s Sitting
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			return nil, fmt.Errorf("mock: sitting log line %d: %w", len(out)+1, err)
		}
		out = append(out, s)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("mock: read sitting log: %w", err)
	}
	return out, nil
}

// HistoryLine summarizes past sittings for the mock start screen.
func HistoryLine(sittings []Sitting) string {
	if len(sittings) == 0 {
		return ""
	}
	best := 0
	for _, s := range sittings {
		if n := s.Solved(); n > best {
			best = n
		}
	}
	last := sittings[len(sittings)-1].Solved()
	noun := "sittings"
	if len(sittings) == 1 {
		noun = "sitting"
	}
	return fmt.Sprintf("%d %s · best %d/4 · last %d/4", len(sittings), noun, best, last)
}
