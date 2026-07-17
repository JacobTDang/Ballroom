package main

import "testing"

func TestMatchTable(t *testing.T) {
	cases := []struct {
		pattern, s string
		want       bool
	}{
		// literals & ?
		{"abc", "abc", true},
		{"abc", "abd", false},
		{"a?c", "abc", true},
		{"a?c", "ac", false},
		{"?", "", false},
		// stars
		{"*", "", true},
		{"*", "anything", true},
		{"a*", "a", true},
		{"a*b*c", "aXXbYYc", true},
		{"a*b*c", "aXXbYY", false},
		{"*.go", "main.go", true},
		{"*.go", "main.gox", false},
		{"a*a", "aa", true},
		{"a*a", "aba", true},
		{"a*a", "ab", false},
		{"**", "x", true},
		// classes
		{"[a-c]x", "bx", true},
		{"[a-c]x", "dx", false},
		{"[xyz]", "y", true},
		{"[xyz]", "w", false},
		{"file[0-9].txt", "file7.txt", true},
		{"file[0-9].txt", "fileX.txt", false},
		// unclosed class: invalid pattern matches nothing
		{"[abc", "a", false},
		{"[", "[", false},
	}
	for _, c := range cases {
		if got := Match(c.pattern, c.s); got != c.want {
			t.Errorf("Match(%q, %q) = %v, want %v", c.pattern, c.s, got, c.want)
		}
	}
}
