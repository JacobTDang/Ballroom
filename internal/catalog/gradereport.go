package catalog

import (
	"regexp"
	"sort"
	"strings"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// DimensionRating is one rubric dimension's rating parsed out of a
// design grader's summary (tracker.Attempt.GradeSummary).
type DimensionRating struct {
	Name   string
	Rating string // "strong" | "adequate" | "missing", lower-cased
}

// DimensionWeakness aggregates one dimension's ratings across attempts
// -- the data behind the Stats screen's "rubric weak spots" section.
type DimensionWeakness struct {
	Name     string
	Missing  int
	Adequate int
	Strong   int
}

// Total is how many graded attempts rated this dimension at all.
func (d DimensionWeakness) Total() int {
	return d.Missing + d.Adequate + d.Strong
}

// dimensionRatingPattern matches the grading contract's per-dimension
// lines tolerantly: optional markdown bold, a leading case number, the
// dimension name, then the rating word. Live grader outputs seen so
// far: "1. Use cases & constraints: missing. ..." (hy3) and
// "**1. Use cases & constraints**: Adequate. ..." (laguna).
var dimensionRatingPattern = regexp.MustCompile(`(?i)^\s*(?:\*\*)?\s*\d+\.\s*([^:*]+?)\s*(?:\*\*)?\s*:\s*(?:\*\*)?\s*(strong|adequate|missing)\b`)

// ParseDimensionRatings extracts per-dimension ratings from a grading
// summary. Freeform prose that doesn't follow the dimension-line shape
// parses to nothing -- absence of signal, not an error, since grader
// formatting is model behavior and Stats must degrade gracefully.
func ParseDimensionRatings(summary string) []DimensionRating {
	var out []DimensionRating
	for _, line := range strings.Split(summary, "\n") {
		m := dimensionRatingPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		out = append(out, DimensionRating{
			Name:   strings.TrimSpace(m[1]),
			Rating: strings.ToLower(m[2]),
		})
	}
	return out
}

// WeakDimensions aggregates dimension ratings across every attempt that
// carries a grade summary, worst dimensions first (most missing, then
// most adequate). Attempts without summaries -- coding attempts,
// self-assessed design attempts -- contribute nothing.
func WeakDimensions(attempts []tracker.Attempt) []DimensionWeakness {
	byKey := map[string]*DimensionWeakness{}
	var order []string
	for _, a := range attempts {
		if a.GradeSummary == "" {
			continue
		}
		for _, r := range ParseDimensionRatings(a.GradeSummary) {
			key := strings.ToLower(r.Name)
			d, ok := byKey[key]
			if !ok {
				d = &DimensionWeakness{Name: r.Name}
				byKey[key] = d
				order = append(order, key)
			}
			switch r.Rating {
			case "missing":
				d.Missing++
			case "adequate":
				d.Adequate++
			case "strong":
				d.Strong++
			}
		}
	}
	out := make([]DimensionWeakness, 0, len(order))
	for _, key := range order {
		out = append(out, *byKey[key])
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Missing != out[j].Missing {
			return out[i].Missing > out[j].Missing
		}
		return out[i].Adequate > out[j].Adequate
	})
	return out
}
