package catalog

import (
	"testing"
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
