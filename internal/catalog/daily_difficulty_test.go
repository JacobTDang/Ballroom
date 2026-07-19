package catalog

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func dailyProblem(id, difficulty string, solved bool) ProblemStatus {
	v := ExerciseStatus{Exercise: exercise.Exercise{
		ID: id + "-go", ProblemID: id, Title: id, Category: exercise.CategoryArraysHashing,
		Language: exercise.LanguageGo, Difficulty: difficulty,
	}}
	if solved {
		v.Attempts = 1
		v.LastResult = tracker.ResultPass
		v.LastAttemptDate = time.Now().Format("2006-01-02")
	}
	return ProblemStatus{
		ProblemID: id, Title: id, Category: exercise.CategoryArraysHashing,
		Solved: solved, Variants: []ExerciseStatus{v},
	}
}

func mixedDifficultyCatalog(solvedCount int) []ProblemStatus {
	var out []ProblemStatus
	for i := 0; i < solvedCount; i++ {
		out = append(out, dailyProblem("solved-"+string(rune('a'+i%26))+string(rune('0'+i/26)), exercise.DifficultyEasy, true))
	}
	for i := 0; i < 8; i++ {
		out = append(out, dailyProblem("easy-"+string(rune('a'+i)), exercise.DifficultyEasy, false))
		out = append(out, dailyProblem("med-"+string(rune('a'+i)), exercise.DifficultyMedium, false))
		out = append(out, dailyProblem("hard-"+string(rune('a'+i)), exercise.DifficultyHard, false))
	}
	return out
}

func pickDifficulty(p ProblemStatus) string {
	for _, v := range p.Variants {
		if v.Exercise.Difficulty != "" {
			return v.Exercise.Difficulty
		}
	}
	return ""
}

// TestDailyPick_BeginnerOnlyGetsEasy is the whole point: an unweighted
// hash over ~600 unsolved problems can hand someone on day one an
// advanced graph problem, which is a bad first assignment even though
// it's a valid one.
func TestDailyPick_BeginnerOnlyGetsEasy(t *testing.T) {
	problems := mixedDifficultyCatalog(0)
	// Sweep dates rather than trusting one: the pick is date-derived,
	// so a single date proves nothing about the gate.
	for d := 0; d < 60; d++ {
		day := time.Now().AddDate(0, 0, d)
		got, ok := DailyPick(problems, day)
		if !ok {
			t.Fatalf("day %d: no pick", d)
		}
		if diff := pickDifficulty(got); diff != exercise.DifficultyEasy {
			t.Fatalf("day %d: picked %s problem %q for a profile with nothing solved", d, diff, got.ProblemID)
		}
	}
}

func TestDailyPick_WidensToMediumWithProgress(t *testing.T) {
	problems := mixedDifficultyCatalog(20) // past the easy-only threshold
	sawMedium := false
	for d := 0; d < 90; d++ {
		got, ok := DailyPick(problems, time.Now().AddDate(0, 0, d))
		if !ok {
			t.Fatal("no pick")
		}
		if diff := pickDifficulty(got); diff == exercise.DifficultyHard {
			t.Fatalf("day %d: hard problem %q at 20 solved, want the pool still capped at medium", d, got.ProblemID)
		} else if diff == exercise.DifficultyMedium {
			sawMedium = true
		}
	}
	if !sawMedium {
		t.Error("never picked a medium problem at 20 solved -- the pool did not widen")
	}
}

func TestDailyPick_AllDifficultiesOnceExperienced(t *testing.T) {
	problems := mixedDifficultyCatalog(70)
	sawHard := false
	for d := 0; d < 90; d++ {
		got, _ := DailyPick(problems, time.Now().AddDate(0, 0, d))
		if pickDifficulty(got) == exercise.DifficultyHard {
			sawHard = true
			break
		}
	}
	if !sawHard {
		t.Error("never picked a hard problem at 70 solved -- the pool never fully opened")
	}
}

// TestDailyPick_StableWithinADay preserves the original contract: the
// pick is an assignment for the day, not a slot machine.
func TestDailyPick_StableWithinADay(t *testing.T) {
	problems := mixedDifficultyCatalog(0)
	day := time.Now()
	first, _ := DailyPick(problems, day)
	for i := 0; i < 20; i++ {
		got, _ := DailyPick(problems, day)
		if got.ProblemID != first.ProblemID {
			t.Fatalf("pick changed within the same day: %q then %q", first.ProblemID, got.ProblemID)
		}
	}
}

func TestDailyPick_ChangesTomorrow(t *testing.T) {
	problems := mixedDifficultyCatalog(0)
	today, _ := DailyPick(problems, time.Now())
	differs := false
	for d := 1; d < 10; d++ {
		if got, _ := DailyPick(problems, time.Now().AddDate(0, 0, d)); got.ProblemID != today.ProblemID {
			differs = true
			break
		}
	}
	if !differs {
		t.Error("pick never changed across ten days")
	}
}

// TestDailyPick_DueWorkBeatsTheDifficultyGate: due problems are
// time-sensitive, and gating them away would silently drop review work
// a beginner has already started.
func TestDailyPick_DueWorkBeatsTheDifficultyGate(t *testing.T) {
	hardDue := dailyProblem("hard-due", exercise.DifficultyHard, false)
	hardDue.Variants[0].Attempts = 1
	hardDue.Variants[0].LastResult = tracker.ResultFail
	hardDue.Variants[0].LastAttemptDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	hardDue.Attempts = 1

	problems := append(mixedDifficultyCatalog(0), hardDue)
	found := false
	for d := 0; d < 60; d++ {
		if got, _ := DailyPick(problems, time.Now().AddDate(0, 0, d)); got.ProblemID == "hard-due" {
			found = true
			break
		}
	}
	if !found {
		t.Error("a due hard problem was never picked -- the gate swallowed review work")
	}
}

// TestDailyPick_WidensRatherThanReturningNothing: if the gated pool is
// empty (all easy problems solved, still under the threshold), the pick
// must widen instead of leaving the user with no assignment.
func TestDailyPick_WidensRatherThanReturningNothing(t *testing.T) {
	var problems []ProblemStatus
	for i := 0; i < 5; i++ {
		problems = append(problems, dailyProblem("easy-"+string(rune('a'+i)), exercise.DifficultyEasy, true))
	}
	problems = append(problems, dailyProblem("hard-only", exercise.DifficultyHard, false))

	got, ok := DailyPick(problems, time.Now())
	if !ok {
		t.Fatal("no pick when only hard problems remain unsolved")
	}
	if got.ProblemID != "hard-only" {
		t.Errorf("picked %q, want the only unsolved problem", got.ProblemID)
	}
}

func TestDailyPick_EmptyCatalogStillReturnsFalse(t *testing.T) {
	if _, ok := DailyPick(nil, time.Now()); ok {
		t.Error("empty catalog returned a pick")
	}
}

// TestDailyPick_WidensOneRankAtATime: when the gated pool is empty,
// jumping straight to "anything unsolved" can hand a beginner a hard
// problem the gate existed to withhold. Stepping up one rank keeps the
// progression intact -- a user with every easy problem done gets a
// medium, not a hard.
func TestDailyPick_WidensOneRankAtATime(t *testing.T) {
	var problems []ProblemStatus
	for i := 0; i < 5; i++ { // under the 15-solved threshold
		problems = append(problems, dailyProblem("easy-"+string(rune('a'+i)), exercise.DifficultyEasy, true))
	}
	problems = append(problems, dailyProblem("med-left", exercise.DifficultyMedium, false))
	problems = append(problems, dailyProblem("hard-left", exercise.DifficultyHard, false))

	for d := 0; d < 30; d++ {
		got, ok := DailyPick(problems, time.Now().AddDate(0, 0, d))
		if !ok {
			t.Fatalf("day %d: no pick", d)
		}
		if got.ProblemID == "hard-left" {
			t.Fatalf("day %d: widened past medium to a hard problem while a medium was available", d)
		}
	}
}
