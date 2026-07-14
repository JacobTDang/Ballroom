package main

import "testing"

func TestCharacterReplacement(t *testing.T) {
	cases := []struct {
		s    string
		k    int
		want int
	}{
		{"ABAB", 2, 4},
		{"AABABBA", 1, 4},
		{"ABCDE", 1, 2},
		{"AAAA", 0, 4},
		{"A", 0, 1},
		{"ABBB", 2, 4},
		{"", 2, 0},
		{"AAAA", 4, 4},
		{"ABABABAB", 3, 7},
	}

	for _, c := range cases {
		got := CharacterReplacement(c.s, c.k)
		if got != c.want {
			t.Errorf("CharacterReplacement(%q, %d) = %d, want %d", c.s, c.k, got, c.want)
		}
	}
}
