package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Tier 1: exact small-n correctness.
func TestBuildLog_SmallCases(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"c", "b", "a"}, "abc"},
		{[]string{"single"}, "single"},
		{[]string{}, ""},
		{[]string{"d", "c", "b", "a"}, "abcd"},
		{[]string{"world", "hello"}, "helloworld"},
		{[]string{"x", "x", "x"}, "xxx"},
	}
	for _, c := range cases {
		if got := BuildLog(c.in); got != c.want {
			t.Errorf("BuildLog(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Tier 2: an in-test stopwatch, not a harness timeout. n=200,000
// chunks of 15 bytes each (~3MB total copying for a correct,
// linear-time build):
//
//	fixed:     ~3,000,000 bytes copied total -- well under 100ms even
//	           unoptimized -- >=100x headroom under the 10s bound.
//	quadratic: ~15 * 200,000^2 / 2 =~ 3*10^11 bytes (~300GB) of
//	           copying -- tens of seconds (20-40s observed) -- >=2x
//	           OVER the 10s bound.
//
// The 10s cutoff sits comfortably between those two, so this cannot
// flake based on machine speed; it only distinguishes O(n) from
// O(n^2). time.Now()/time.Since carry Go's monotonic clock reading, so
// this isn't affected by wall-clock adjustments either.
func TestBuildLog_LargeInputFinishesWithinTenSeconds(t *testing.T) {
	const n = 200_000
	chunks := make([]string, n)
	for i := range chunks {
		chunks[i] = fmt.Sprintf("%014d;", i) // 15 bytes each
	}

	start := time.Now()
	result := BuildLog(chunks)
	elapsed := time.Since(start)

	if len(result) != n*15 {
		t.Fatalf("BuildLog result length = %d, want %d", len(result), n*15)
	}
	if !strings.HasPrefix(result, chunks[n-1]) {
		t.Errorf("BuildLog result doesn't start with the oldest chunk")
	}
	if !strings.HasSuffix(result, chunks[0]) {
		t.Errorf("BuildLog result doesn't end with the newest chunk")
	}
	if elapsed >= 10*time.Second {
		t.Errorf("BuildLog took %s on %d chunks -- looks quadratic", elapsed, n)
	}
}
