package main

import "testing"

func TestMinWindow(t *testing.T) {
	cases := []struct {
		s, tt string
		want  string
	}{
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"a", "a", "a"},
		{"a", "aa", ""},
		{"ab", "b", "b"},
		{"bba", "ab", "ba"},
		{"abc", "abc", "abc"},
		{"aaflslflsldkalskaaa", "aaa", "aaa"},
		{"cabwefgewcwaefgcf", "cae", "cwae"},
	}

	for _, c := range cases {
		got := MinWindow(c.s, c.tt)
		if got != c.want {
			t.Errorf("MinWindow(%q, %q) = %q, want %q", c.s, c.tt, got, c.want)
		}
	}
}
