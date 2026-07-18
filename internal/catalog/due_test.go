package catalog

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func namedProblem(id string) ProblemStatus {
	return ProblemStatus{ProblemID: id}
}

func TestSortDueFirst_FloatsDueProblemsKeepingRelativeOrder(t *testing.T) {
	due := map[string]bool{"c": true, "e": true}
	problems := []ProblemStatus{
		namedProblem("a"), namedProblem("b"), namedProblem("c"), namedProblem("d"), namedProblem("e"),
	}

	got := SortDueFirst(problems, func(p ProblemStatus) bool { return due[p.ProblemID] })

	want := []string{"c", "e", "a", "b", "d"}
	for i, w := range want {
		if got[i].ProblemID != w {
			t.Fatalf("order = %v, want %v -- due problems first, both partitions keeping their original relative order",
				problemIDs(got), want)
		}
	}
}

func TestSortDueFirst_DoesNotMutateItsInput(t *testing.T) {
	problems := []ProblemStatus{namedProblem("a"), namedProblem("b")}

	SortDueFirst(problems, func(p ProblemStatus) bool { return p.ProblemID == "b" })

	if problems[0].ProblemID != "a" || problems[1].ProblemID != "b" {
		t.Errorf("input mutated to %v -- callers hold this slice as appModel.problems", problemIDs(problems))
	}
}

func TestSortDueFirst_NoDueProblemsIsIdentity(t *testing.T) {
	problems := []ProblemStatus{namedProblem("b"), namedProblem("a"), namedProblem("c")}

	got := SortDueFirst(problems, func(ProblemStatus) bool { return false })

	want := []string{"b", "a", "c"}
	for i, w := range want {
		if got[i].ProblemID != w {
			t.Fatalf("order = %v, want the input order %v unchanged", problemIDs(got), want)
		}
	}
}

func problemIDs(problems []ProblemStatus) []string {
	ids := make([]string, len(problems))
	for i, p := range problems {
		ids[i] = p.ProblemID
	}
	return ids
}

func attemptedVariant(lang, lastDate, lastResult string) ExerciseStatus {
	return ExerciseStatus{
		Exercise:        exercise.Exercise{Kind: exercise.KindDesign, Language: lang},
		Attempts:        1,
		LastResult:      lastResult,
		LastAttemptDate: lastDate,
	}
}

func TestReviewDue_Table(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		p    ProblemStatus
		want bool
	}{
		{
			"never attempted is new, not review",
			ProblemStatus{Variants: []ExerciseStatus{{Exercise: exercise.Exercise{Language: "go"}}}},
			false,
		},
		{
			"failed 3 days ago is due for a retry",
			ProblemStatus{Variants: []ExerciseStatus{attemptedVariant("go", "2026-07-13", tracker.ResultFail)}},
			true,
		},
		{
			"failed yesterday is too fresh",
			ProblemStatus{Variants: []ExerciseStatus{attemptedVariant("go", "2026-07-15", tracker.ResultFail)}},
			false,
		},
		{
			// issue #238: giving up (Solved stays false -- see
			// problems_test.go's GaveUpDoesNotCountAsSolved) resurfaces on
			// the same unsolved 3-day cadence a fail does, not the
			// solved-only 30-day one.
			"gave up 3 days ago is due for a retry",
			ProblemStatus{Variants: []ExerciseStatus{attemptedVariant("go", "2026-07-13", tracker.ResultGaveUp)}},
			true,
		},
		{
			"solved 30 days ago is due for a refresh",
			ProblemStatus{Solved: true, Variants: []ExerciseStatus{attemptedVariant("go", "2026-06-16", tracker.ResultPass)}},
			true,
		},
		{
			"solved a week ago is not",
			ProblemStatus{Solved: true, Variants: []ExerciseStatus{attemptedVariant("go", "2026-07-09", tracker.ResultPass)}},
			false,
		},
		{
			"solved problem uses the newest attempt across variants",
			ProblemStatus{Solved: true, Variants: []ExerciseStatus{
				attemptedVariant("go", "2026-05-01", tracker.ResultPass),
				attemptedVariant("py", "2026-07-10", tracker.ResultFail),
			}},
			false, // touched 6 days ago, even though the pass itself is old
		},
		{
			"unparseable date degrades to not due, never panics",
			ProblemStatus{Variants: []ExerciseStatus{attemptedVariant("go", "not-a-date", tracker.ResultFail)}},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ReviewDue(c.p, now); got != c.want {
				t.Errorf("ReviewDue = %v, want %v", got, c.want)
			}
		})
	}
}

func TestDue_UnionOfMockAndReview(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)

	mockDue := ProblemStatus{Variants: []ExerciseStatus{
		{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass, LastAttemptDate: "2026-07-15"},
		{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
	}}
	reviewDue := ProblemStatus{Variants: []ExerciseStatus{attemptedVariant("go", "2026-07-01", tracker.ResultFail)}}
	neither := ProblemStatus{Variants: []ExerciseStatus{{Exercise: exercise.Exercise{Language: "go"}}}}

	if !Due(mockDue, now) {
		t.Error("Due = false for a mock-due problem, want true")
	}
	if !Due(reviewDue, now) {
		t.Error("Due = false for a review-due problem, want true")
	}
	if Due(neither, now) {
		t.Error("Due = true for an untouched problem, want false")
	}
}
