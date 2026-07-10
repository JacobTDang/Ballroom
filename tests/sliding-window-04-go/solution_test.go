package main

import "testing"

func TestCheckInclusion(t *testing.T) {
	cases := []struct {
		s1, s2 string
		want   bool
	}{
		{"ab", "eidbaooo", true},
		{"ab", "eidboaoo", false},
		{"adc", "dcda", true},
		{"hello", "ooolleoooleh", false},
		{"a", "a", true},
		{"abc", "ab", false},
	}

	for _, c := range cases {
		got := CheckInclusion(c.s1, c.s2)
		if got != c.want {
			t.Errorf("CheckInclusion(%q, %q) = %v, want %v", c.s1, c.s2, got, c.want)
		}
	}
}
