package catalog

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

const hy3Summary = `VERDICT: fail
1. Use cases & constraints: missing. The design only says...
2. Back-of-envelope estimates: missing. No assumed volume appears.
3. High-level design: adequate. Both flows present but thin.
4. Short-code generation: strong. Base-62 with length justified.`

const boldSummary = `VERDICT: pass

**1. Use cases & constraints**: Adequate. Scope stated upfront.
**2. Back-of-envelope estimates**: Strong. Arithmetic shown.`

func TestParseDimensionRatings_PlainNumberedFormat(t *testing.T) {
	got := ParseDimensionRatings(hy3Summary)
	if len(got) != 4 {
		t.Fatalf("parsed %d dimensions, want 4: %+v", len(got), got)
	}
	if got[0].Name != "Use cases & constraints" || got[0].Rating != "missing" {
		t.Errorf("first dimension = %+v, want Use cases & constraints/missing", got[0])
	}
	if got[3].Name != "Short-code generation" || got[3].Rating != "strong" {
		t.Errorf("last dimension = %+v, want Short-code generation/strong", got[3])
	}
}

func TestParseDimensionRatings_BoldMarkdownFormat(t *testing.T) {
	got := ParseDimensionRatings(boldSummary)
	if len(got) != 2 {
		t.Fatalf("parsed %d dimensions, want 2: %+v", len(got), got)
	}
	if got[0].Rating != "adequate" || got[1].Rating != "strong" {
		t.Errorf("ratings = %q,%q, want adequate,strong (case-normalized)", got[0].Rating, got[1].Rating)
	}
}

func TestParseDimensionRatings_ProseWithoutRatingsYieldsNothing(t *testing.T) {
	if got := ParseDimensionRatings("Looks good overall, nice work on the estimates."); len(got) != 0 {
		t.Errorf("parsed %+v from unstructured prose, want none", got)
	}
}

func TestWeakDimensions_AggregatesAcrossAttemptsWorstFirst(t *testing.T) {
	attempts := []tracker.Attempt{
		{Category: "system-design", GradeSummary: hy3Summary},
		{Category: "system-design", GradeSummary: "1. Back-of-envelope estimates: missing. Again no numbers."},
		{Category: "system-design", GradeSummary: ""},  // self-assessed: no signal
		{Category: "arrays-hashing", GradeSummary: ""}, // coding: ignored
	}
	got := WeakDimensions(attempts)
	if len(got) == 0 {
		t.Fatal("WeakDimensions returned nothing")
	}
	if got[0].Name != "Back-of-envelope estimates" || got[0].Missing != 2 {
		t.Errorf("worst dimension = %+v, want Back-of-envelope estimates with 2 missing", got[0])
	}
	for _, d := range got {
		if d.Name == "Short-code generation" && d.Missing != 0 {
			t.Errorf("Short-code generation Missing = %d, want 0", d.Missing)
		}
	}
}
