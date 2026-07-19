package catalog

import (
	"fmt"
	"sort"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// Recommend answers "what should I work on next" -- the question 645
// exercises and eight roadmaps make genuinely hard to answer alone. The
// progress bars show where you stand but suggest nothing, and Daily is
// one date-stable pick; this is the deliberate version.
//
// It reads the same signals the rest of the app already computes
// (MockDue, ReviewDue, difficulty, rubric weak spots) rather than
// inventing a curriculum format -- a roadmap file the app has to parse
// is a maintenance surface, and these signals are already true.

// RecommendKind distinguishes why something is being suggested, so the
// UI can order and phrase them.
type RecommendKind int

const (
	// RecommendDue is time-sensitive: a stale failure, a solved problem
	// gone cold, or a coach pass whose mock is untouched.
	RecommendDue RecommendKind = iota
	// RecommendNext is forward progress in the track you've done least.
	RecommendNext
	// RecommendWeakSpot targets a rubric dimension you keep losing.
	RecommendWeakSpot
)

// Recommendation is one suggestion with the reason it was made. The
// reason is shown verbatim -- a suggestion the user can't evaluate is
// just another thing to distrust.
type Recommendation struct {
	Problem ProblemStatus
	Kind    RecommendKind
	Reason  string
}

const maxRecommendations = 3

// Recommend returns up to three suggestions, strongest signal first.
func Recommend(problems []ProblemStatus, attempts []tracker.Attempt, now time.Time) []Recommendation {
	if len(problems) == 0 {
		return nil
	}

	var out []Recommendation
	taken := map[string]bool{}
	add := func(r Recommendation) {
		if len(out) >= maxRecommendations || taken[r.Problem.ProblemID] {
			return
		}
		taken[r.Problem.ProblemID] = true
		out = append(out, r)
	}

	if due, ok := oldestDue(problems, now); ok {
		add(due)
	}
	if next, ok := nextUnsolved(problems); ok {
		add(next)
	}
	if weak, ok := weakSpotDrill(problems, attempts, taken); ok {
		add(weak)
	}

	// Everything solved and nothing due yet: rather than an empty
	// screen, offer the least-recently-touched solved problem as a
	// refresher.
	if len(out) == 0 {
		if stale, ok := stalestSolved(problems); ok {
			add(Recommendation{
				Problem: stale,
				Kind:    RecommendDue,
				Reason:  "everything's solved — this one's gone coldest",
			})
		}
	}
	return out
}

// oldestDue prefers whatever has been waiting longest -- a mock that
// never happened or a failure gone stale.
func oldestDue(problems []ProblemStatus, now time.Time) (Recommendation, bool) {
	type candidate struct {
		p      ProblemStatus
		reason string
		when   string // last attempt date, "" sorts first (never touched)
	}
	var cands []candidate
	for _, p := range problems {
		switch {
		case MockDue(p):
			cands = append(cands, candidate{p, "coach pass done — the timed mock is still untouched", problemLastDate(p)})
		case ReviewDue(p, now):
			reason := "a failure worth another pass"
			if p.Solved {
				reason = "solved a while back — worth checking it still comes out clean"
			}
			cands = append(cands, candidate{p, reason, problemLastDate(p)})
		}
	}
	if len(cands) == 0 {
		return Recommendation{}, false
	}
	sort.SliceStable(cands, func(a, b int) bool { return cands[a].when < cands[b].when })
	c := cands[0]
	return Recommendation{Problem: c.p, Kind: RecommendDue, Reason: c.reason}, true
}

// nextUnsolved picks forward progress in the least-progressed track,
// gated to a difficulty the user is ready for -- the same gate Daily
// applies, for the same reason: handing a beginner a hard graph problem
// is technically "next" and practically useless.
func nextUnsolved(problems []ProblemStatus) (Recommendation, bool) {
	allowed := allowedDifficulties(solvedCount(problems))

	solvedByTrack := map[string]int{}
	totalByTrack := map[string]int{}
	for _, p := range problems {
		track := TopLevelGroup(p.Category)
		totalByTrack[track]++
		if p.Solved {
			solvedByTrack[track]++
		}
	}

	// Ties are the common case early on (every track sits at zero), and
	// Go randomizes map iteration -- ranging over totalByTrack directly
	// made the suggestion change on every refresh, which reads as the
	// app being indecisive. Walk the tracks in the catalog's own order
	// so a tie always resolves the same way.
	tracks := make([]string, 0, len(totalByTrack))
	for track := range totalByTrack {
		tracks = append(tracks, track)
	}
	sort.SliceStable(tracks, func(a, b int) bool { return CategoryRank(tracks[a]) < CategoryRank(tracks[b]) })

	best := ""
	bestRatio := 2.0
	for _, track := range tracks {
		total := totalByTrack[track]
		if total == 0 || solvedByTrack[track] >= total {
			continue
		}
		if ratio := float64(solvedByTrack[track]) / float64(total); ratio < bestRatio {
			bestRatio, best = ratio, track
		}
	}
	if best == "" {
		return Recommendation{}, false
	}

	pick, ok := firstUnsolvedIn(problems, best, allowed)
	if !ok {
		// Nothing at the gated difficulty left in that track; take any
		// unsolved problem rather than skipping the slot.
		if pick, ok = firstUnsolvedIn(problems, best, nil); !ok {
			return Recommendation{}, false
		}
	}
	return Recommendation{
		Problem: pick,
		Kind:    RecommendNext,
		Reason:  fmt.Sprintf("next in %s — your least-practiced track", DisplayCategory(best)),
	}, true
}

func firstUnsolvedIn(problems []ProblemStatus, track string, allowed map[string]bool) (ProblemStatus, bool) {
	for _, p := range problems {
		if p.Solved || TopLevelGroup(p.Category) != track {
			continue
		}
		if allowed != nil && !allowed[problemDifficulty(p)] {
			continue
		}
		return p, true
	}
	return ProblemStatus{}, false
}

// weakSpotDrill points at a track whose rubric dimension keeps costing
// points. It only fires with real grading history behind it.
func weakSpotDrill(problems []ProblemStatus, attempts []tracker.Attempt, taken map[string]bool) (Recommendation, bool) {
	weak := WeakDimensions(attempts)
	if len(weak) == 0 {
		return Recommendation{}, false
	}
	for _, p := range problems {
		if taken[p.ProblemID] || p.Solved {
			continue
		}
		if p.Category != exercise.CategorySystemDesign && p.Category != exercise.CategoryAPIDesign && p.Category != exercise.CategoryBehavioral {
			continue
		}
		return Recommendation{
			Problem: p,
			Kind:    RecommendWeakSpot,
			Reason:  fmt.Sprintf("practice for %q — the dimension you lose most often", weak[0].Name),
		}, true
	}
	return Recommendation{}, false
}

func stalestSolved(problems []ProblemStatus) (ProblemStatus, bool) {
	var best ProblemStatus
	bestDate := ""
	found := false
	for _, p := range problems {
		if !p.Solved {
			continue
		}
		d := problemLastDate(p)
		if !found || d < bestDate {
			best, bestDate, found = p, d, true
		}
	}
	return best, found
}

func problemLastDate(p ProblemStatus) string {
	latest := ""
	for _, v := range p.Variants {
		if v.LastAttemptDate > latest {
			latest = v.LastAttemptDate
		}
	}
	return latest
}

// problemDifficulty reports the problem's difficulty from its first
// rated variant. Variants of one problem share a difficulty except for
// the design tracks' coach/interviewer split, where the coach variant
// comes first and is the one a recommendation should lead with.
func problemDifficulty(p ProblemStatus) string {
	for _, v := range p.Variants {
		if v.Exercise.Difficulty != "" {
			return v.Exercise.Difficulty
		}
	}
	return ""
}

func solvedCount(problems []ProblemStatus) int {
	n := 0
	for _, p := range problems {
		if p.Solved {
			n++
		}
	}
	return n
}

// difficultyRanks is the easy-to-hard order, used when a gated pool
// comes up empty and has to widen a rank at a time.
var difficultyRanks = []string{exercise.DifficultyEasy, exercise.DifficultyMedium, exercise.DifficultyHard}

// allowedDifficulties widens as you make progress. An unrated problem
// is always allowed -- absent metadata should never hide content.
func allowedDifficulties(solved int) map[string]bool {
	switch {
	case solved < 15:
		return map[string]bool{exercise.DifficultyEasy: true, "": true}
	case solved < 60:
		return map[string]bool{exercise.DifficultyEasy: true, exercise.DifficultyMedium: true, "": true}
	default:
		return nil // everything
	}
}
