package main

import "testing"

func TestFirstAtLeast(t *testing.T) {
	cases := []struct {
		in     []int
		target int
		want   int
	}{
		{[]int{1, 3, 5, 7, 9}, 6, 3},
		{[]int{1, 3, 5, 7, 9}, 1, 0},
		{[]int{1, 3, 5, 7, 9}, 0, 0},
		{[]int{1, 3, 5, 7, 9}, 9, 4},
		{[]int{1, 3, 5, 7, 9}, 10, 5},
		{[]int{5}, 10, 1},
		{[]int{2, 2, 2, 2}, 2, 0},
	}
	for _, c := range cases {
		got := FirstAtLeast(c.in, c.target)
		if got != c.want {
			t.Errorf("FirstAtLeast(%v, %d) = %d, want %d", c.in, c.target, got, c.want)
		}
	}
}
