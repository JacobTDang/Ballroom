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
	var candidates []ProblemStatus
	for _, p := range problems {
		if Due(p, now) || !p.Solved {
			candidates = append(candidates, p)
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
