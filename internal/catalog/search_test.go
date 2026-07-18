package catalog

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func searchFixture() []ProblemStatus {
	return []ProblemStatus{
		{ProblemID: "two-sum-01", Title: "Two Sum", Category: exercise.CategoryArraysHashing,
			Variants: []ExerciseStatus{{Exercise: exercise.Exercise{ID: "two-sum-01-go"}}, {Exercise: exercise.Exercise{ID: "two-sum-01-python"}}}},
		{ProblemID: "two-pointers-01", Title: "Two Sum II", Category: exercise.CategoryTwoPointers,
			Variants: []ExerciseStatus{{Exercise: exercise.Exercise{ID: "two-pointers-01-go"}}}},
		{ProblemID: "api-quotas-01", Title: "Rate Limits & Quotas for the Platform", Category: exercise.CategoryAPIDesign,
			Variants: []ExerciseStatus{{Exercise: exercise.Exercise{ID: "api-quotas-01-coach"}}}},
		{ProblemID: "rate-limiter-01", Title: "Fixed-Window Rate Limiter", Category: exercise.CategoryImplementation,
			Variants: []ExerciseStatus{{Exercise: exercise.Exercise{ID: "rate-limiter-01-go"}}}},
	}
}

func titlesOf(got []ProblemStatus) []string {
	out := make([]string, 0, len(got))
	for _, p := range got {
		out = append(out, p.Title)
	}
	return out
}

func TestSearch_EmptyQueryReturnsEverything(t *testing.T) {
	all := searchFixture()
	if got := Search(all, ""); len(got) != len(all) {
		t.Errorf("Search(\"\") returned %d, want all %d", len(got), len(all))
	}
}

// TestSearch_ExactExerciseIDRanksFirst: typing an id you already know
// should land on it, not bury it under fuzzy title matches.
func TestSearch_ExactExerciseIDRanksFirst(t *testing.T) {
	got := Search(searchFixture(), "two-pointers-01-go")
	if len(got) == 0 {
		t.Fatal("exact exercise id matched nothing")
	}
	if got[0].ProblemID != "two-pointers-01" {
		t.Errorf("first result = %q, want two-pointers-01", got[0].ProblemID)
	}
}

func TestSearch_ProblemIDPrefixRanksAboveFuzzyTitle(t *testing.T) {
	got := Search(searchFixture(), "rate-limiter")
	if len(got) == 0 {
		t.Fatal("problem id prefix matched nothing")
	}
	if got[0].ProblemID != "rate-limiter-01" {
		t.Errorf("first result = %q, want the id-prefix match rate-limiter-01, got order %v", got[0].ProblemID, titlesOf(got))
	}
}

func TestSearch_MatchesTitleCaseInsensitively(t *testing.T) {
	got := Search(searchFixture(), "RATE LIMITS")
	if len(got) != 1 || got[0].ProblemID != "api-quotas-01" {
		t.Errorf("got %v, want just the Rate Limits & Quotas question", titlesOf(got))
	}
}

// TestSearch_MatchesCategoryName is what makes "show me everything in
// api design" work without knowing any titles.
func TestSearch_MatchesCategoryName(t *testing.T) {
	got := Search(searchFixture(), "api design")
	if len(got) != 1 || got[0].ProblemID != "api-quotas-01" {
		t.Errorf("got %v, want the api-design problem", titlesOf(got))
	}
}

// TestSearch_FuzzySubsequence: "twsm" should still find "Two Sum" —
// this is the affordance that makes typing fast.
func TestSearch_FuzzySubsequence(t *testing.T) {
	got := Search(searchFixture(), "twsm")
	if len(got) == 0 {
		t.Fatal("fuzzy subsequence matched nothing")
	}
	found := false
	for _, p := range got {
		if p.ProblemID == "two-sum-01" {
			found = true
		}
	}
	if !found {
		t.Errorf("got %v, want Two Sum among the fuzzy matches", titlesOf(got))
	}
}

func TestSearch_NoMatchReturnsEmpty(t *testing.T) {
	if got := Search(searchFixture(), "zzzzqqqq"); len(got) != 0 {
		t.Errorf("got %v, want no matches", titlesOf(got))
	}
}

// TestSearch_StableForTies keeps the result order predictable so the
// cursor doesn't jump around between keystrokes that don't change the
// match set's ranking.
func TestSearch_StableForTies(t *testing.T) {
	first := Search(searchFixture(), "two")
	second := Search(searchFixture(), "two")
	if len(first) != len(second) {
		t.Fatalf("unstable result count: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].ProblemID != second[i].ProblemID {
			t.Errorf("unstable order at %d: %q vs %q", i, first[i].ProblemID, second[i].ProblemID)
		}
	}
}
