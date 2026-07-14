package main

import "testing"

func TestContainsDuplicate(t *testing.T) {
	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3, 4}, false},
		{[]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{[]int{1}, false},
		{[]int{1, 1}, true},
		{[]int{1, 2}, false},
		{[]int{-1, -1}, true},
		{[]int{-5, -3, -1, 1, 3, 5}, false},
		{[]int{0, 4, 5, 0, 3, 6}, true},
		{[]int{7, 7, 7, 7, 7}, true},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1}, true},
		{[]int{-1000000000, 1000000000}, false},
		{[]int{1000000000, 1000000000}, true},
	}
	for _, c := range cases {
		got := ContainsDuplicate(c.nums)
		if got != c.want {
			t.Errorf("ContainsDuplicate(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}
