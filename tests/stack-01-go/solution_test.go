package main

import "testing"

func TestIsValid(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
		{"([)]", false},
		{"{[]}", true},
		{"(", false},
		{"]", false},
	}

	for _, c := range cases {
		got := IsValid(c.s)
		if got != c.want {
			t.Errorf("IsValid(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}
