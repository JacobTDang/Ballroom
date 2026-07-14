package main

import "testing"

func TestIsAnagram(t *testing.T) {
	cases := []struct {
		s, t string
		want bool
	}{
		{"anagram", "nagaram", true},
		{"rat", "car", false},
		{"ab", "a", false},
		{"aacc", "ccac", false},
		{"a", "a", true},
		{"aabbcc", "abcabc", true},
		{"listen", "silent", true},
		{"aaab", "aabb", false},
		{"a", "b", false},
		{"abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcba", true},
		{"aaaa", "aaaa", true},
	}
	for _, c := range cases {
		got := IsAnagram(c.s, c.t)
		if got != c.want {
			t.Errorf("IsAnagram(%q, %q) = %v, want %v", c.s, c.t, got, c.want)
		}
	}
}
