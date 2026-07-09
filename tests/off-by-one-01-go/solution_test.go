package main

import "testing"

func TestMaxOf(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{3, 1, 4, 1, 5, 9, 2, 6}, 9},
		{[]int{-5, -1, -10}, -1},
		{[]int{42}, 42},
		{[]int{5, 5, 5}, 5},              // tie for max
		{[]int{1, 2, 3, 4, 5, 100}, 100}, // max is the last element (the exact boundary the bug touches)
		{[]int{-1, -2, -3, -100}, -1},    // max is the first element, all negative
	}
	for _, c := range cases {
		got := MaxOf(c.in)
		if got != c.want {
			t.Errorf("MaxOf(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}
