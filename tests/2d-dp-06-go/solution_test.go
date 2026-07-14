package main

import "testing"

func TestIsInterleave_Classic(t *testing.T) {
	if got := IsInterleave("aabcc", "dbbca", "aadbbcbcac"); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "aabcc", "dbbca", "aadbbcbcac", got)
	}
}

func TestIsInterleave_NotInterleaved(t *testing.T) {
	if got := IsInterleave("aabcc", "dbbca", "aadbbbaccc"); got != false {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want false", "aabcc", "dbbca", "aadbbbaccc", got)
	}
}

func TestIsInterleave_AllEmpty(t *testing.T) {
	if got := IsInterleave("", "", ""); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "", "", "", got)
	}
}

func TestIsInterleave_OneEmpty(t *testing.T) {
	if got := IsInterleave("a", "", "a"); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "a", "", "a", got)
	}
}

func TestIsInterleave_LengthMismatch(t *testing.T) {
	if got := IsInterleave("abc", "def", "abcde"); got != false {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want false", "abc", "def", "abcde", got)
	}
}

func TestIsInterleave_FirstEmptyMatch(t *testing.T) {
	if got := IsInterleave("", "abc", "abc"); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "", "abc", "abc", got)
	}
}

func TestIsInterleave_FirstEmptyMismatch(t *testing.T) {
	if got := IsInterleave("", "abc", "abd"); got != false {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want false", "", "abc", "abd", got)
	}
}

func TestIsInterleave_AmbiguousMultipleWays(t *testing.T) {
	if got := IsInterleave("ab", "ab", "abab"); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "ab", "ab", "abab", got)
	}
}

func TestIsInterleave_RequiresBacktrackChoice(t *testing.T) {
	if got := IsInterleave("ab", "ab", "aabb"); got != true {
		t.Errorf("IsInterleave(%q, %q, %q) = %v, want true", "ab", "ab", "aabb", got)
	}
}
