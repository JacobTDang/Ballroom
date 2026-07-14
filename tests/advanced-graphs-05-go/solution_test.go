package main

import "testing"

// isValidAlienOrder checks the ORDERING PROPERTY rather than an exact
// string, since the topological order implied by words is not unique:
// every distinct character appearing in words must appear exactly once
// in order, and every adjacent word pair's first differing character
// must respect that order.
func isValidAlienOrder(words []string, order string) bool {
	if order == "" {
		return false
	}

	pos := make(map[byte]int)
	for i := 0; i < len(order); i++ {
		if _, exists := pos[order[i]]; exists {
			return false // duplicate character in order
		}
		pos[order[i]] = i
	}

	seen := make(map[byte]bool)
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			seen[w[i]] = true
		}
	}
	if len(seen) != len(pos) {
		return false
	}
	for c := range seen {
		if _, ok := pos[c]; !ok {
			return false
		}
	}

	for i := 0; i < len(words)-1; i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		if len(w1) > len(w2) && w1[:minLen] == w2[:minLen] {
			return false
		}
		for j := 0; j < minLen; j++ {
			if w1[j] != w2[j] {
				if pos[w1[j]] >= pos[w2[j]] {
					return false
				}
				break
			}
		}
	}
	return true
}

func TestAlienOrder_Valid(t *testing.T) {
	words := []string{"wrt", "wrf", "er", "ett", "rftt"}
	order := AlienOrder(words)
	if !isValidAlienOrder(words, order) {
		t.Errorf("AlienOrder(%v) = %q, not a valid ordering", words, order)
	}
}

func TestAlienOrder_InvalidPrefix(t *testing.T) {
	words := []string{"abc", "ab"}
	if got := AlienOrder(words); got != "" {
		t.Errorf("AlienOrder(%v) = %q, want empty (invalid prefix order)", words, got)
	}
}

func TestAlienOrder_Cycle(t *testing.T) {
	words := []string{"z", "x", "z"}
	if got := AlienOrder(words); got != "" {
		t.Errorf("AlienOrder(%v) = %q, want empty (cycle)", words, got)
	}
}

func TestAlienOrder_SingleWord(t *testing.T) {
	if got := AlienOrder([]string{"z"}); got != "z" {
		t.Errorf(`AlienOrder(["z"]) = %q, want "z"`, got)
	}
}

func TestAlienOrder_TwoDistinctChars(t *testing.T) {
	if got := AlienOrder([]string{"a", "b"}); got != "ab" {
		t.Errorf(`AlienOrder(["a", "b"]) = %q, want "ab"`, got)
	}
}
