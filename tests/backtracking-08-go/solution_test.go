package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestLetterCombinations(t *testing.T) {
	cases := []struct {
		digits string
		want   []string
	}{
		{"23", []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}},
		{"", []string{}},
		{"2", []string{"a", "b", "c"}},
		{"9", []string{"w", "x", "y", "z"}},
		{
			"79",
			[]string{"pw", "px", "py", "pz", "qw", "qx", "qy", "qz", "rw", "rx", "ry", "rz", "sw", "sx", "sy", "sz"},
		},
	}

	for _, c := range cases {
		got := append([]string(nil), LetterCombinations(c.digits)...)
		sort.Strings(got)
		want := append([]string(nil), c.want...)
		sort.Strings(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("LetterCombinations(%q) = %v, want %v (order-independent)", c.digits, got, want)
		}
	}
}
