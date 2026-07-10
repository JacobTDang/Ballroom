package main

import "testing"

func TestLadderLength_Classic(t *testing.T) {
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}
	if got := LadderLength("hit", "cog", wordList); got != 5 {
		t.Errorf("LadderLength(hit, cog, %v) = %d, want 5", wordList, got)
	}
}

func TestLadderLength_EndWordNotInList(t *testing.T) {
	wordList := []string{"hot", "dot", "dog", "lot", "log"}
	if got := LadderLength("hit", "cog", wordList); got != 0 {
		t.Errorf("LadderLength(hit, cog, %v) = %d, want 0", wordList, got)
	}
}

func TestLadderLength_DirectNeighbor(t *testing.T) {
	wordList := []string{"hot"}
	if got := LadderLength("hit", "hot", wordList); got != 2 {
		t.Errorf("LadderLength(hit, hot, %v) = %d, want 2", wordList, got)
	}
}
