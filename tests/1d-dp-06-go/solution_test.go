package main

import "testing"

func TestCountSubstrings_ThreeDistinct(t *testing.T) {
	if got := CountSubstrings("abc"); got != 3 {
		t.Errorf("CountSubstrings(%q) = %d, want 3", "abc", got)
	}
}

func TestCountSubstrings_AllSame(t *testing.T) {
	if got := CountSubstrings("aaa"); got != 6 {
		t.Errorf("CountSubstrings(%q) = %d, want 6", "aaa", got)
	}
}

func TestCountSubstrings_OddPalindrome(t *testing.T) {
	if got := CountSubstrings("aba"); got != 4 {
		t.Errorf("CountSubstrings(%q) = %d, want 4", "aba", got)
	}
}

func TestCountSubstrings_SingleChar(t *testing.T) {
	if got := CountSubstrings("z"); got != 1 {
		t.Errorf("CountSubstrings(%q) = %d, want 1", "z", got)
	}
}

func TestCountSubstrings_TwoSame(t *testing.T) {
	if got := CountSubstrings("aa"); got != 3 {
		t.Errorf("CountSubstrings(%q) = %d, want 3", "aa", got)
	}
}

func TestCountSubstrings_TwoDifferent(t *testing.T) {
	if got := CountSubstrings("ab"); got != 2 {
		t.Errorf("CountSubstrings(%q) = %d, want 2", "ab", got)
	}
}

func TestCountSubstrings_LargerAllSame(t *testing.T) {
	if got := CountSubstrings("aaaaa"); got != 15 {
		t.Errorf("CountSubstrings(%q) = %d, want 15", "aaaaa", got)
	}
}

func TestCountSubstrings_NestedPalindromes(t *testing.T) {
	if got := CountSubstrings("aabaa"); got != 9 {
		t.Errorf("CountSubstrings(%q) = %d, want 9", "aabaa", got)
	}
}
