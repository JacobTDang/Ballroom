package main

import "testing"

func TestWordBreak_Classic(t *testing.T) {
	wordDict := []string{"leet", "code"}
	if got := WordBreak("leetcode", wordDict); got != true {
		t.Errorf("WordBreak(%q, %v) = %v, want true", "leetcode", wordDict, got)
	}
}

func TestWordBreak_ReusedWord(t *testing.T) {
	wordDict := []string{"apple", "pen"}
	if got := WordBreak("applepenapple", wordDict); got != true {
		t.Errorf("WordBreak(%q, %v) = %v, want true", "applepenapple", wordDict, got)
	}
}

func TestWordBreak_Impossible(t *testing.T) {
	wordDict := []string{"cats", "dog", "sand", "and", "cat"}
	if got := WordBreak("catsandog", wordDict); got != false {
		t.Errorf("WordBreak(%q, %v) = %v, want false", "catsandog", wordDict, got)
	}
}

func TestWordBreak_SingleWord(t *testing.T) {
	wordDict := []string{"a"}
	if got := WordBreak("a", wordDict); got != true {
		t.Errorf("WordBreak(%q, %v) = %v, want true", "a", wordDict, got)
	}
}

func TestWordBreak_LeftoverCharUnmatched(t *testing.T) {
	wordDict := []string{"a"}
	if got := WordBreak("ab", wordDict); got != false {
		t.Errorf("WordBreak(%q, %v) = %v, want false", "ab", wordDict, got)
	}
}

func TestWordBreak_TrailingCharNeverMatches(t *testing.T) {
	wordDict := []string{"a", "aa"}
	if got := WordBreak("aaaaaaaaaaaaaaaaaaaab", wordDict); got != false {
		t.Errorf("WordBreak(%q, %v) = %v, want false", "aaaaaaaaaaaaaaaaaaaab", wordDict, got)
	}
}

func TestWordBreak_MultiplePaths(t *testing.T) {
	wordDict := []string{"car", "ca", "rs"}
	if got := WordBreak("cars", wordDict); got != true {
		t.Errorf("WordBreak(%q, %v) = %v, want true", "cars", wordDict, got)
	}
}

func TestWordBreak_SimpleConcatenation(t *testing.T) {
	wordDict := []string{"goal", "special"}
	if got := WordBreak("goalspecial", wordDict); got != true {
		t.Errorf("WordBreak(%q, %v) = %v, want true", "goalspecial", wordDict, got)
	}
}
