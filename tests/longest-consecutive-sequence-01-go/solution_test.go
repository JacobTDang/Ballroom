package main

import "testing"

func TestLongestConsecutive(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{100, 4, 200, 1, 3, 2}, 4},
		{[]int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, 9},
		{[]int{}, 0},
		{[]int{1, 2, 0, 1}, 3},
		{[]int{9, 1, 4, 7, 3, -1, 0, 5, 8, -1, 6}, 7},
		{[]int{5}, 1},
		{[]int{7, 7, 7, 7}, 1},
		{[]int{1, 2, 3, 10, 11}, 3},
		{[]int{-3, -2, -1, 0, 1}, 5},
		{[]int{50, 3, 51, 2, 52, 1, 4, 49, 48, 47}, 6},
		{[]int{-1000000000, -999999999, -999999998}, 3},
	}
	for _, c := range cases {
		got := LongestConsecutive(c.nums)
		if got != c.want {
			t.Errorf("LongestConsecutive(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}
