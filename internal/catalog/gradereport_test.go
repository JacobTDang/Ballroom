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

func TestCodingWeakSpots_RanksByFailRatioWithThreshold(t *testing.T) {
	attempts := []tracker.Attempt{
		// trees: 3 attempts, 2 fails -- worst.
		{Category: "trees", Result: tracker.ResultFail},
		{Category: "trees", Result: tracker.ResultFail},
		{Category: "trees", Result: tracker.ResultPass},
		// 1d-dp: 4 attempts, 1 fail -- weak but better.
		{Category: "1d-dp", Result: tracker.ResultFail},
		{Category: "1d-dp", Result: tracker.ResultPass},
		{Category: "1d-dp", Result: tracker.ResultPass},
		{Category: "1d-dp", Result: tracker.ResultPass},
		// graphs: only 2 attempts -- under the threshold, excluded even
		// though both failed.
		{Category: "graphs", Result: tracker.ResultFail},
		{Category: "graphs", Result: tracker.ResultFail},
		// heap: 3 attempts, all passing -- not a weak spot at all.
		{Category: "heap", Result: tracker.ResultPass},
		{Category: "heap", Result: tracker.ResultPass},
		{Category: "heap", Result: tracker.ResultPass},
	}
	got := CodingWeakSpots(attempts, 3)
	if len(got) != 2 {
		t.Fatalf("CodingWeakSpots = %+v, want exactly trees and 1d-dp", got)
	}
	if got[0].Category != "trees" || got[0].Fails != 2 || got[0].Attempts != 3 {
		t.Errorf("worst spot = %+v, want trees 2/3", got[0])
	}
	if got[1].Category != "1d-dp" || got[1].Fails != 1 || got[1].Attempts != 4 {
		t.Errorf("second spot = %+v, want 1d-dp 1/4", got[1])
	}
}

func TestCodingWeakSpots_ExcludesDesignAndBehavioralTracks(t *testing.T) {
	attempts := []tracker.Attempt{
		{Category: "system-design", Result: tracker.ResultFail},
		{Category: "system-design", Result: tracker.ResultFail},
		{Category: "system-design", Result: tracker.ResultFail},
		{Category: "behavioral", Result: tracker.ResultFail},
		{Category: "behavioral", Result: tracker.ResultFail},
		{Category: "behavioral", Result: tracker.ResultFail},
	}
	if got := CodingWeakSpots(attempts, 3); len(got) != 0 {
		t.Errorf("CodingWeakSpots = %+v, want design/behavioral attempts excluded (they have the rubric section)", got)
	}
}

func TestCodingWeakSpots_TieBreaksByMoreAttempts(t *testing.T) {
	attempts := []tracker.Attempt{
		// Both at 50% fail ratio; stack has more attempts, so more
		// evidence -- it ranks first.
		{Category: "greedy", Result: tracker.ResultFail},
		{Category: "greedy", Result: tracker.ResultFail},
		{Category: "greedy", Result: tracker.ResultPass},
		{Category: "greedy", Result: tracker.ResultPass},
		{Category: "stack", Result: tracker.ResultFail},
		{Category: "stack", Result: tracker.ResultFail},
		{Category: "stack", Result: tracker.ResultFail},
		{Category: "stack", Result: tracker.ResultPass},
		{Category: "stack", Result: tracker.ResultPass},
		{Category: "stack", Result: tracker.ResultPass},
	}
	got := CodingWeakSpots(attempts, 3)
	if len(got) != 2 || got[0].Category != "stack" || got[1].Category != "greedy" {
		t.Errorf("CodingWeakSpots = %+v, want stack (6 attempts) ranked above greedy (4) on equal ratio", got)
	}
}

func TestCodingWeakSpots_EmptyAndUnderThresholdDegradeToNothing(t *testing.T) {
	if got := CodingWeakSpots(nil, 3); len(got) != 0 {
		t.Errorf("CodingWeakSpots(nil) = %+v, want empty", got)
	}
	one := []tracker.Attempt{{Category: "tries", Result: tracker.ResultFail}}
	if got := CodingWeakSpots(one, 3); len(got) != 0 {
		t.Errorf("one attempt under threshold = %+v, want empty", got)
	}
}
