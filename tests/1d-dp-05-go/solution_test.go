package main

import "testing"

func TestLongestPalindrome_OddTie(t *testing.T) {
	if got := LongestPalindrome("babad"); got != "bab" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "babad", got, "bab")
	}
}

func TestLongestPalindrome_Even(t *testing.T) {
	if got := LongestPalindrome("cbbd"); got != "bb" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "cbbd", got, "bb")
	}
}

func TestLongestPalindrome_SingleChar(t *testing.T) {
	if got := LongestPalindrome("a"); got != "a" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "a", got, "a")
	}
}

func TestLongestPalindrome_WholeString(t *testing.T) {
	if got := LongestPalindrome("abba"); got != "abba" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "abba", got, "abba")
	}
}

func TestLongestPalindrome_AllSameLonger(t *testing.T) {
	if got := LongestPalindrome("aaaaa"); got != "aaaaa" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "aaaaa", got, "aaaaa")
	}
}

func TestLongestPalindrome_NoRepeat(t *testing.T) {
	if got := LongestPalindrome("abcde"); got != "a" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "abcde", got, "a")
	}
}

func TestLongestPalindrome_BuriedInLargerString(t *testing.T) {
	if got := LongestPalindrome("zzabcbayy"); got != "abcba" {
		t.Errorf("LongestPalindrome(%q) = %q, want %q", "zzabcbayy", got, "abcba")
	}
}
