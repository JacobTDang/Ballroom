package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{" ", true},
		{"0P", false},
		{"Was it a car or a cat I saw?", true},
		{".,", true},
		{"a_b", false},
		{"12321", true},
		{"ab", false},
		{"", true},
		{"Able was I, ere I saw Elba", true},
	}

	for _, c := range cases {
		got := IsPalindrome(c.s)
		if got != c.want {
			t.Errorf("IsPalindrome(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}
