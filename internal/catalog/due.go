package catalog

import (
	"sort"
	"time"
)

// Review-due thresholds: a failed problem resurfaces quickly (the gap
// is for a bit of forgetting, not abandonment); a solved one rests a
// month before a refresh pass.
const (
	reviewRetryAfter   = 3 * 24 * time.Hour
	reviewRefreshAfter = 30 * 24 * time.Hour
)

// ReviewDue is the date-based half of "due" (MockDue is the
// progression half): a problem whose last attempt failed at least 3
// days ago is due for a retry, and a solved problem untouched for 30
// days is due for a refresh. Freshness is the newest attempt across
// all variants -- practicing the Python variant yesterday means the
// problem was touched yesterday, whatever the Go variant's state.
// Never-attempted problems are new, not "review"; an unparseable date
// (hand-edited DB rows) degrades to not-due rather than panicking a
// render loop.
func ReviewDue(p ProblemStatus, now time.Time) bool {
	last := ""
	for _, v := range p.Variants {
		// "2006-01-02" sorts lexicographically, so string max is date max.
		if v.LastAttemptDate > last {
			last = v.LastAttemptDate
		}
	}
	if last == "" {
		return false
	}
	d, err := time.Parse("2006-01-02", last)
	if err != nil {
		return false
	}
	age := now.Sub(d)
	if p.Solved {
		return age >= reviewRefreshAfter
	}
	return age >= reviewRetryAfter
}

// Due is what the picker sorts and marks by: due for the roadmap's
// interviewer second pass (MockDue) or due by date (ReviewDue).
func Due(p ProblemStatus, now time.Time) bool {
	return MockDue(p) || ReviewDue(p, now)
}

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
