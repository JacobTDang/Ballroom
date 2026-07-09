package catalog

import "github.com/JacobTDang/Ballroom/internal/tracker"

// ProblemStatus groups an exercise's language variants (same ProblemID)
// into one problem-level view — this is what the practice picker
// actually shows and lets you select, with the specific language chosen
// afterward.
type ProblemStatus struct {
	ProblemID string
	Title     string
	Category  string
	Variants  []ExerciseStatus // one per available language
	Solved    bool             // true if any variant's most recent attempt passed
	Attempts  int              // summed across all variants
}

// GroupByProblem groups a flat, already-loaded exercise list by
// ProblemID, preserving each problem's first-encountered order (so
// callers that pre-sort statuses, like List, keep that ordering).
// Solved is language-agnostic: a problem is solved if any of its
// variants has been solved, regardless of which language.
func GroupByProblem(statuses []ExerciseStatus) []ProblemStatus {
	var order []string
	byProblem := make(map[string]*ProblemStatus)

	for _, s := range statuses {
		pid := s.Exercise.ProblemID
		p, ok := byProblem[pid]
		if !ok {
			p = &ProblemStatus{
				ProblemID: pid,
				Title:     s.Exercise.Title,
				Category:  s.Exercise.Category,
			}
			byProblem[pid] = p
			order = append(order, pid)
		}
		p.Variants = append(p.Variants, s)
		p.Attempts += s.Attempts
		if s.LastResult == tracker.ResultPass {
			p.Solved = true
		}
	}

	problems := make([]ProblemStatus, len(order))
	for i, pid := range order {
		problems[i] = *byProblem[pid]
	}
	return problems
}
