package main

import "testing"

func TestLongestCommonSubsequence_Classic(t *testing.T) {
	if got := LongestCommonSubsequence("abcde", "ace"); got != 3 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 3", "abcde", "ace", got)
	}
}

func TestLongestCommonSubsequence_Identical(t *testing.T) {
	if got := LongestCommonSubsequence("abc", "abc"); got != 3 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 3", "abc", "abc", got)
	}
}

func TestLongestCommonSubsequence_NoCommon(t *testing.T) {
	if got := LongestCommonSubsequence("abc", "def"); got != 0 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 0", "abc", "def", got)
	}
}

func TestLongestCommonSubsequence_EmptyFirst(t *testing.T) {
	if got := LongestCommonSubsequence("", "abc"); got != 0 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 0", "", "abc", got)
	}
}

func TestLongestCommonSubsequence_DifferentOrder(t *testing.T) {
	if got := LongestCommonSubsequence("abc", "acb"); got != 2 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 2", "abc", "acb", got)
	}
}

func TestLongestCommonSubsequence_InterspersedNoise(t *testing.T) {
	if got := LongestCommonSubsequence("aggtab", "gxtxayb"); got != 4 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 4", "aggtab", "gxtxayb", got)
	}
}

func TestLongestCommonSubsequence_RepeatedChars(t *testing.T) {
	if got := LongestCommonSubsequence("aaaa", "aa"); got != 2 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 2", "aaaa", "aa", got)
	}
}

func TestLongestCommonSubsequence_SingleCharMatch(t *testing.T) {
	if got := LongestCommonSubsequence("a", "a"); got != 1 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 1", "a", "a", got)
	}
}

func TestLongestCommonSubsequence_BothEmpty(t *testing.T) {
	if got := LongestCommonSubsequence("", ""); got != 0 {
		t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want 0", "", "", got)
	}
}
