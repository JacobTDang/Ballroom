package main

import "testing"

func TestLargestRectangleArea(t *testing.T) {
	cases := []struct {
		heights []int
		want    int
	}{
		{[]int{2, 1, 5, 6, 2, 3}, 10},
		{[]int{2, 4}, 4},
		{[]int{1}, 1},
		{[]int{0, 0}, 0},
		{[]int{5, 5, 5, 5}, 20},
		{[]int{5, 4, 3, 2, 1}, 9},
		{[]int{1, 2, 3, 4, 5}, 9},
	}

	for _, c := range cases {
		got := LargestRectangleArea(c.heights)
		if got != c.want {
			t.Errorf("LargestRectangleArea(%v) = %d, want %d", c.heights, got, c.want)
		}
	}
}
