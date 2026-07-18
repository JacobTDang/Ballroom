package catalog

import (
	"hash/fnv"
	"time"
)

// DailyPick is the main menu's "Daily" problem: a deterministic pick
// for the date -- same problem all day, a different one tomorrow -- so
// it reads as an assignment rather than a slot machine. Candidates are
// problems that are due (catalog.Due) or never solved; when everything
// is solved and fresh, any problem qualifies (a refresh day). ok is
// false only for an empty catalog.
//
// FNV-1a over the date string, mod the candidate count: no state to
// store, and every session that day agrees on the pick.
func DailyPick(problems []ProblemStatus, now time.Time) (ProblemStatus, bool) {
	// Difficulty gates the pool by how far along you are. An unweighted
	// hash over ~600 unsolved problems will eventually hand someone on
	// day one an advanced graph problem -- a valid pick and a bad
	// assignment. Due work is exempt: it's time-sensitive, and gating it
	// would silently drop review a beginner has already started.
	allowed := allowedDifficulties(solvedCount(problems))

	var candidates []ProblemStatus
	for _, p := range problems {
		if Due(p, now) {
			candidates = append(candidates, p)
			continue
		}
		if !p.Solved && (allowed == nil || allowed[problemDifficulty(p)]) {
			candidates = append(candidates, p)
		}
	}
	// Nothing at the gated difficulty (every easy problem solved while
	// still under the next threshold): widen one rank at a time rather
	// than all the way open. Jumping straight to "anything unsolved"
	// would hand out exactly the hard problem the gate existed to
	// withhold, when a medium was sitting right there.
	for rank := 0; len(candidates) == 0 && rank < len(difficultyRanks); rank++ {
		widened := map[string]bool{"": true}
		for i := 0; i <= rank; i++ {
			widened[difficultyRanks[i]] = true
		}
		for _, p := range problems {
			if !p.Solved && widened[problemDifficulty(p)] {
				candidates = append(candidates, p)
			}
		}
	}
	if len(candidates) == 0 {
		candidates = problems
	}
	if len(candidates) == 0 {
		return ProblemStatus{}, false
	}
	h := fnv.New32a()
	h.Write([]byte(now.Format("2006-01-02")))
	return candidates[int(h.Sum32()%uint32(len(candidates)))], true
}
