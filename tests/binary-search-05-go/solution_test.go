package main

import "testing"

func TestSearch(t *testing.T) {
	cases := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{[]int{1}, 0, -1},
		{[]int{5, 1, 3}, 5, 0},
		{[]int{1, 3}, 3, 1},
		{[]int{9, 10, 1, 2, 3, 4, 5, 6, 7, 8}, 8, 9},
		{[]int{9, 10, 1, 2, 3, 4, 5, 6, 7, 8}, 100, -1},
	}

	for _, c := range cases {
		got := Search(c.nums, c.target)
		if got != c.want {
			t.Errorf("Search(%v, %d) = %d, want %d", c.nums, c.target, got, c.want)
		}
	}
}
