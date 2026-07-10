package main

import "testing"

func TestSearch(t *testing.T) {
	cases := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{[]int{-1, 0, 3, 5, 9, 12}, 2, -1},
		{[]int{5}, 5, 0},
		{[]int{2, 5}, 5, 1},
		{[]int{2, 5}, 1, -1},
	}

	for _, c := range cases {
		got := Search(c.nums, c.target)
		if got != c.want {
			t.Errorf("Search(%v, %d) = %d, want %d", c.nums, c.target, got, c.want)
		}
	}
}
