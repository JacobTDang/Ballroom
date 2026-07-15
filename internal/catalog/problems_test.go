package catalog

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func fakeVariant(problemID, id, category, language, title string, attempts int, lastResult string) ExerciseStatus {
	return ExerciseStatus{
		Exercise: exercise.Exercise{
			ID:        id,
			ProblemID: problemID,
			Title:     title,
			Category:  category,
			Language:  language,
		},
		Attempts:   attempts,
		LastResult: lastResult,
	}
}

func TestGroupByProblem_GroupsLanguageVariantsTogether(t *testing.T) {
	statuses := []ExerciseStatus{
		fakeVariant("two-pointers-01", "two-pointers-01-go", "pattern", "go", "Two Sum II", 0, ""),
		fakeVariant("two-pointers-01", "two-pointers-01-cpp", "pattern", "cpp", "Two Sum II", 0, ""),
		fakeVariant("two-pointers-01", "two-pointers-01-python", "pattern", "python", "Two Sum II", 0, ""),
	}

	problems := GroupByProblem(statuses)
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem grouping 3 variants, got %d", len(problems))
	}
	if len(problems[0].Variants) != 3 {
		t.Errorf("expected 3 variants, got %d", len(problems[0].Variants))
	}
	if problems[0].ProblemID != "two-pointers-01" {
		t.Errorf("ProblemID = %q, want %q", problems[0].ProblemID, "two-pointers-01")
	}
	if problems[0].Title != "Two Sum II" {
		t.Errorf("Title = %q, want %q", problems[0].Title, "Two Sum II")
	}
	if problems[0].Category != "pattern" {
		t.Errorf("Category = %q, want %q", problems[0].Category, "pattern")
	}
}

func TestGroupByProblem_SolvedIfAnyVariantPassed(t *testing.T) {
	statuses := []ExerciseStatus{
		fakeVariant("p1", "p1-go", "pattern", "go", "P1", 2, tracker.ResultFail),
		fakeVariant("p1", "p1-cpp", "pattern", "cpp", "P1", 1, tracker.ResultPass),
		fakeVariant("p1", "p1-python", "pattern", "python", "P1", 0, ""),
	}

	problems := GroupByProblem(statuses)
	if !problems[0].Solved {
		t.Error("expected Solved=true since the cpp variant passed, even though go failed and python was untouched")
	}
}

func TestGroupByProblem_NotSolvedIfNoVariantPassed(t *testing.T) {
	statuses := []ExerciseStatus{
		fakeVariant("p1", "p1-go", "pattern", "go", "P1", 1, tracker.ResultFail),
		fakeVariant("p1", "p1-cpp", "pattern", "cpp", "P1", 0, ""),
	}

	problems := GroupByProblem(statuses)
	if problems[0].Solved {
		t.Error("expected Solved=false when no variant has passed")
	}
}

func TestGroupByProblem_AttemptsSummedAcrossVariants(t *testing.T) {
	statuses := []ExerciseStatus{
		fakeVariant("p1", "p1-go", "pattern", "go", "P1", 2, tracker.ResultFail),
		fakeVariant("p1", "p1-cpp", "pattern", "cpp", "P1", 3, tracker.ResultPass),
	}

	problems := GroupByProblem(statuses)
	if problems[0].Attempts != 5 {
		t.Errorf("Attempts = %d, want 5 (summed across variants)", problems[0].Attempts)
	}
}

func TestGroupByProblem_MultipleProblemsPreserveInputOrder(t *testing.T) {
	statuses := []ExerciseStatus{
		fakeVariant("two-pointers-01", "two-pointers-01-go", "pattern", "go", "Two Sum II", 0, ""),
		fakeVariant("off-by-one-01", "off-by-one-01-go", "debug", "go", "Off by one", 0, ""),
		fakeVariant("two-pointers-01", "two-pointers-01-cpp", "pattern", "cpp", "Two Sum II", 0, ""),
	}

	problems := GroupByProblem(statuses)
	if len(problems) != 2 {
		t.Fatalf("expected 2 distinct problems, got %d", len(problems))
	}
	// first-encountered order, matching List()'s pre-sorted input (by
	// category then id) rather than re-sorting.
	if problems[0].ProblemID != "two-pointers-01" || problems[1].ProblemID != "off-by-one-01" {
		t.Errorf("expected order [two-pointers-01, off-by-one-01], got [%s, %s]",
			problems[0].ProblemID, problems[1].ProblemID)
	}
	if len(problems[0].Variants) != 2 {
		t.Errorf("expected two-pointers-01's two variants grouped despite the debug problem in between, got %d",
			len(problems[0].Variants))
	}
}

func TestMockDue_CoachPassedInterviewerUnattempted(t *testing.T) {
	p := ProblemStatus{
		ProblemID: "url-shortener-01",
		Category:  exercise.CategorySystemDesign,
		Variants: []ExerciseStatus{
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass},
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}, Attempts: 0},
		},
	}
	if !MockDue(p) {
		t.Error("MockDue = false, want true: coach passed, interviewer never attempted")
	}
}

func TestMockDue_FalseCases(t *testing.T) {
	base := func() ProblemStatus {
		return ProblemStatus{
			ProblemID: "url-shortener-01",
			Category:  exercise.CategorySystemDesign,
			Variants: []ExerciseStatus{
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass},
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}, Attempts: 0},
			},
		}
	}

	p := base() // coach never passed
	p.Variants[0].LastResult = tracker.ResultFail
	if MockDue(p) {
		t.Error("MockDue = true with coach failed, want false")
	}

	p = base() // interviewer already attempted
	p.Variants[1].Attempts = 2
	if MockDue(p) {
		t.Error("MockDue = true with interviewer attempted, want false")
	}

	p = base() // interviewer-only mock: no coach variant at all
	p.Variants = p.Variants[1:]
	if MockDue(p) {
		t.Error("MockDue = true for an interviewer-only problem, want false")
	}

	p = base() // coding problem shape never qualifies
	for i := range p.Variants {
		p.Variants[i].Exercise.Kind = exercise.KindCoding
	}
	if MockDue(p) {
		t.Error("MockDue = true for a coding problem, want false")
	}
}
