package catalog

import (
	"sort"
	"strings"
)

// Search finds problems across the whole catalog, not just one
// category — the picker's own filter only ever sees the category you
// already drilled into, which means finding one of 645 problems used to
// require remembering its category first. Matching covers what a person
// actually remembers: a title fragment, a problem or exercise id, or the
// track's name.
//
// Ranking runs strongest-signal-first (exact id, then id prefix, then
// title prefix, then substring, then a fuzzy subsequence) because the
// alternative — one flat "contains" pass — buries the id you typed in
// full under every problem whose title happens to share a word.
func Search(problems []ProblemStatus, query string) []ProblemStatus {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return problems
	}

	type scored struct {
		problem ProblemStatus
		rank    int
		order   int // input position, so ties never reshuffle between keystrokes
	}
	var hits []scored
	for i, p := range problems {
		if r, ok := matchRank(p, q); ok {
			hits = append(hits, scored{problem: p, rank: r, order: i})
		}
	}
	sort.SliceStable(hits, func(a, b int) bool {
		if hits[a].rank != hits[b].rank {
			return hits[a].rank < hits[b].rank
		}
		return hits[a].order < hits[b].order
	})

	out := make([]ProblemStatus, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.problem)
	}
	return out
}

// Rank values, lowest wins.
const (
	rankExactID = iota
	rankIDPrefix
	rankTitlePrefix
	rankTitleSubstring
	rankCategory
	rankFuzzy
)

func matchRank(p ProblemStatus, q string) (int, bool) {
	problemID := strings.ToLower(p.ProblemID)
	title := strings.ToLower(p.Title)

	if problemID == q {
		return rankExactID, true
	}
	for _, v := range p.Variants {
		if strings.ToLower(v.Exercise.ID) == q {
			return rankExactID, true
		}
	}
	if strings.HasPrefix(problemID, q) {
		return rankIDPrefix, true
	}
	for _, v := range p.Variants {
		if strings.HasPrefix(strings.ToLower(v.Exercise.ID), q) {
			return rankIDPrefix, true
		}
	}
	if strings.HasPrefix(title, q) {
		return rankTitlePrefix, true
	}
	if strings.Contains(title, q) {
		return rankTitleSubstring, true
	}
	// The display name, not the raw id: someone types "api design", not
	// "api-design".
	if strings.Contains(strings.ToLower(DisplayCategory(p.Category)), q) {
		return rankCategory, true
	}
	if subsequence(title, q) {
		return rankFuzzy, true
	}
	return 0, false
}

// subsequence reports whether every rune of q appears in s in order,
// which is what makes "twsm" find "Two Sum". Spaces in the query are
// ignored so a half-typed multi-word title still matches.
func subsequence(s, q string) bool {
	sr := []rune(s)
	i := 0
	for _, want := range q {
		if want == ' ' {
			continue
		}
		for i < len(sr) && sr[i] != want {
			i++
		}
		if i == len(sr) {
			return false
		}
		i++
	}
	return true
}
