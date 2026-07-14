package main

import "testing"

func TestFindKthLargest(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{3, 2, 1, 5, 6, 4}, 2, 5},
		{[]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
		{[]int{1}, 1, 1},
		{[]int{2, 1}, 2, 1},
		{[]int{5, 5, 5, 5}, 2, 5},
		{[]int{-1, -5, -3, -2, -4}, 3, -3},
	}

	for _, c := range cases {
		got := FindKthLargest(c.nums, c.k)
		if got != c.want {
			t.Errorf("FindKthLargest(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
		}
	}
}
