package catalog

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func recProblem(id, title, category, difficulty string, solved bool, attempts int) ProblemStatus {
	return ProblemStatus{
		ProblemID: id, Title: title, Category: category, Solved: solved, Attempts: attempts,
		Variants: []ExerciseStatus{{Exercise: exercise.Exercise{
			ID: id + "-go", ProblemID: id, Title: title, Category: category,
			Language: exercise.LanguageGo, Difficulty: difficulty,
		}}},
	}
}

func recDesign(id, title string, coachSolved, interviewerAttempted bool) ProblemStatus {
	coach := ExerciseStatus{Exercise: exercise.Exercise{
		ID: id + "-coach", ProblemID: id, Kind: exercise.KindDesign,
		Category: exercise.CategorySystemDesign, Language: exercise.LanguageCoach,
	}}
	if coachSolved {
		coach.Attempts = 1
		coach.LastResult = tracker.ResultPass
	}
	interviewer := ExerciseStatus{Exercise: exercise.Exercise{
		ID: id + "-interviewer", ProblemID: id, Kind: exercise.KindDesign,
		Category: exercise.CategorySystemDesign, Language: exercise.LanguageInterviewer,
	}}
	if interviewerAttempted {
		interviewer.Attempts = 1
	}
	return ProblemStatus{
		ProblemID: id, Title: title, Category: exercise.CategorySystemDesign,
		Variants: []ExerciseStatus{coach, interviewer},
	}
}

// TestRecommend_FreshProfileStillSuggestsSomething: the empty-history
// case is the one a new user sees first, so it must never come back
// blank.
func TestRecommend_FreshProfileStillSuggestsSomething(t *testing.T) {
	problems := []ProblemStatus{
		recProblem("two-sum-01", "Two Sum", exercise.CategoryArraysHashing, exercise.DifficultyEasy, false, 0),
		recProblem("word-ladder-01", "Word Ladder", exercise.CategoryGraphs, exercise.DifficultyHard, false, 0),
	}
	got := Recommend(problems, nil, time.Now())
	if len(got) == 0 {
		t.Fatal("Recommend returned nothing for a fresh profile")
	}
	for _, r := range got {
		if r.Reason == "" {
			t.Errorf("recommendation %q has no reason", r.Problem.ProblemID)
		}
		if r.Problem.ProblemID == "" {
			t.Error("recommendation has no problem attached")
		}
	}
}

// TestRecommend_LeadsWithDueWork: due work is time-sensitive in a way
// "next unsolved" isn't, so it has to outrank it.
func TestRecommend_LeadsWithDueWork(t *testing.T) {
	old := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	stale := recProblem("stale-01", "Stale Failure", exercise.CategoryTrees, exercise.DifficultyMedium, false, 1)
	stale.Variants[0].Attempts = 1
	stale.Variants[0].LastResult = tracker.ResultFail
	stale.Variants[0].LastAttemptDate = old

	problems := []ProblemStatus{
		recProblem("fresh-01", "Never Touched", exercise.CategoryArraysHashing, exercise.DifficultyEasy, false, 0),
		stale,
	}
	got := Recommend(problems, nil, time.Now())
	if len(got) == 0 {
		t.Fatal("no recommendations")
	}
	if got[0].Problem.ProblemID != "stale-01" {
		t.Errorf("first recommendation = %q, want the due problem stale-01", got[0].Problem.ProblemID)
	}
	if got[0].Kind != RecommendDue {
		t.Errorf("first recommendation kind = %v, want RecommendDue", got[0].Kind)
	}
}

func TestRecommend_SurfacesAMockDuePass(t *testing.T) {
	problems := []ProblemStatus{
		recDesign("url-shortener-01", "Design Pastebin", true, false),
		recProblem("two-sum-01", "Two Sum", exercise.CategoryArraysHashing, exercise.DifficultyEasy, false, 0),
	}
	got := Recommend(problems, nil, time.Now())
	found := false
	for _, r := range got {
		if r.Problem.ProblemID == "url-shortener-01" && r.Kind == RecommendDue {
			found = true
		}
	}
	if !found {
		t.Errorf("mock-due design problem not recommended: %+v", got)
	}
}

// TestRecommend_NextUnsolvedRespectsDifficultyGate mirrors Daily's
// gating: a beginner shouldn't be pointed at a hard problem.
func TestRecommend_NextUnsolvedRespectsDifficultyGate(t *testing.T) {
	problems := []ProblemStatus{
		recProblem("hard-01", "Hard One", exercise.CategoryGraphs, exercise.DifficultyHard, false, 0),
		recProblem("easy-01", "Easy One", exercise.CategoryGraphs, exercise.DifficultyEasy, false, 0),
	}
	got := Recommend(problems, nil, time.Now())
	for _, r := range got {
		if r.Kind == RecommendNext && r.Problem.ProblemID == "hard-01" {
			t.Errorf("recommended a hard problem to a profile with nothing solved: %+v", r)
		}
	}
}

func TestRecommend_DeduplicatesAcrossKinds(t *testing.T) {
	old := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	p := recProblem("only-01", "The Only Problem", exercise.CategoryTrees, exercise.DifficultyEasy, false, 1)
	p.Variants[0].Attempts = 1
	p.Variants[0].LastResult = tracker.ResultFail
	p.Variants[0].LastAttemptDate = old

	got := Recommend([]ProblemStatus{p}, nil, time.Now())
	seen := map[string]bool{}
	for _, r := range got {
		if seen[r.Problem.ProblemID] {
			t.Errorf("problem %q recommended twice: %+v", r.Problem.ProblemID, got)
		}
		seen[r.Problem.ProblemID] = true
	}
}

func TestRecommend_CapsAtThree(t *testing.T) {
	var problems []ProblemStatus
	for _, id := range []string{"a", "b", "c", "d", "e", "f"} {
		problems = append(problems, recProblem(id+"-01", "Problem "+id, exercise.CategoryTrees, exercise.DifficultyEasy, false, 0))
	}
	if got := Recommend(problems, nil, time.Now()); len(got) > 3 {
		t.Errorf("returned %d recommendations, want at most 3", len(got))
	}
}

func TestRecommend_EmptyCatalogReturnsNothing(t *testing.T) {
	if got := Recommend(nil, nil, time.Now()); len(got) != 0 {
		t.Errorf("got %+v, want nothing for an empty catalog", got)
	}
}

// TestRecommend_AllSolvedFallsBackToReview: once everything is solved
// the "next unsolved" slot has nothing to offer, and the screen must
// not simply go blank.
func TestRecommend_AllSolvedFallsBackToReview(t *testing.T) {
	p := recProblem("done-01", "Finished", exercise.CategoryTrees, exercise.DifficultyEasy, true, 1)
	p.Variants[0].Attempts = 1
	p.Variants[0].LastResult = tracker.ResultPass
	p.Variants[0].LastAttemptDate = time.Now().AddDate(0, 0, -60).Format("2006-01-02")

	got := Recommend([]ProblemStatus{p}, nil, time.Now())
	if len(got) == 0 {
		t.Error("everything solved but nothing recommended -- the 30-day review should resurface")
	}
}

// TestRecommend_StableAcrossCalls: ties are the common case early on
// (every track at zero), and a suggestion that changes on every refresh
// reads as the app being indecisive. Caught live -- the dashboard and
// the picker disagreed because Go randomizes map iteration.
func TestRecommend_StableAcrossCalls(t *testing.T) {
	problems := []ProblemStatus{
		recProblem("dsa-01", "DSA One", exercise.CategoryArraysHashing, exercise.DifficultyEasy, false, 0),
		recProblem("dbg-01", "Debug One", exercise.CategoryDebug, exercise.DifficultyEasy, false, 0),
		recProblem("con-01", "Concurrency One", exercise.CategoryConcurrency, exercise.DifficultyEasy, false, 0),
		recProblem("imp-01", "Implementation One", exercise.CategoryImplementation, exercise.DifficultyEasy, false, 0),
	}
	first := Recommend(problems, nil, time.Now())
	for i := 0; i < 25; i++ {
		got := Recommend(problems, nil, time.Now())
		if len(got) != len(first) {
			t.Fatalf("call %d returned %d recommendations, want %d", i, len(got), len(first))
		}
		for j := range got {
			if got[j].Problem.ProblemID != first[j].Problem.ProblemID {
				t.Fatalf("call %d differs at %d: %q vs %q", i, j, got[j].Problem.ProblemID, first[j].Problem.ProblemID)
			}
		}
	}
}
