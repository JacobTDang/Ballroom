package catalog

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func TestDailyPick_DeterministicForADateAndVariesAcrossDates(t *testing.T) {
	problems := []ProblemStatus{
		namedProblem("a"), namedProblem("b"), namedProblem("c"),
		namedProblem("d"), namedProblem("e"), namedProblem("f"),
		namedProblem("g"), namedProblem("h"), namedProblem("i"),
	}
	day := time.Date(2026, 7, 16, 9, 0, 0, 0, time.UTC)

	first, ok := DailyPick(problems, day)
	if !ok {
		t.Fatal("DailyPick found nothing among 9 candidates")
	}
	// Same date, any time of day: same pick.
	evening, _ := DailyPick(problems, day.Add(11*time.Hour))
	if evening.ProblemID != first.ProblemID {
		t.Errorf("same date picked %q then %q, want a stable daily pick", first.ProblemID, evening.ProblemID)
	}

	// Adjacent dates: the pick must actually rotate. With 9 candidates a
	// single collision is legitimate (hash mod 9), so scan a week --
	// all 7 days landing on one problem would mean the date isn't
	// really feeding the hash.
	distinct := map[string]bool{}
	for i := 0; i < 7; i++ {
		p, _ := DailyPick(problems, day.AddDate(0, 0, i))
		distinct[p.ProblemID] = true
	}
	if len(distinct) < 2 {
		t.Errorf("7 consecutive days all picked %v -- the date isn't influencing the pick", distinct)
	}
}

func TestDailyPick_PrefersDueOrUnsolvedAndFallsBackToAll(t *testing.T) {
	now := time.Date(2026, 7, 16, 9, 0, 0, 0, time.UTC)

	solved := ProblemStatus{ProblemID: "solved-recently", Solved: true, Variants: []ExerciseStatus{
		attemptedVariant("go", "2026-07-15", tracker.ResultPass),
	}}
	unsolved := namedProblem("never-tried")

	// With an unsolved candidate present, a freshly-solved problem must
	// never be the pick.
	for i := 0; i < 7; i++ {
		p, ok := DailyPick([]ProblemStatus{solved, unsolved}, now.AddDate(0, 0, i))
		if !ok || p.ProblemID != "never-tried" {
			t.Fatalf("day %d picked %q, want the unsolved problem while one exists", i, p.ProblemID)
		}
	}

	// Everything solved and fresh: fall back to picking among all.
	p, ok := DailyPick([]ProblemStatus{solved}, now)
	if !ok || p.ProblemID != "solved-recently" {
		t.Errorf("all-solved catalog picked (%q, %v), want the fallback over an empty result", p.ProblemID, ok)
	}

	if _, ok := DailyPick(nil, now); ok {
		t.Error("DailyPick on an empty catalog claimed success")
	}
}
