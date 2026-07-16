package catalog

import "sort"

// SortDueFirst returns problems reordered so every problem due
// satisfies come before the rest, each partition keeping its original
// relative order (List's category-then-id order). The picker applies
// this with MockDue so a due problem can't hide mid-list behind the
// alphabet -- a real gap found while verifying the mock-due marker:
// "search-kv-store" sat below the fold while its marker went unseen.
// Returns a new slice; the input (held live as appModel.problems) is
// never mutated.
func SortDueFirst(problems []ProblemStatus, due func(ProblemStatus) bool) []ProblemStatus {
	out := make([]ProblemStatus, len(problems))
	copy(out, problems)
	sort.SliceStable(out, func(i, j int) bool {
		return due(out[i]) && !due(out[j])
	})
	return out
}
