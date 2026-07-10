package main

import "testing"

func TestFindDuplicate(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{1, 3, 4, 2, 2}, 2},
		{[]int{3, 1, 3, 4, 2}, 3},
		{[]int{1, 1}, 1},
		{[]int{1, 1, 2}, 1},
		{[]int{2, 2, 2, 2, 2}, 2},
	}

	for _, c := range cases {
		got := FindDuplicate(c.nums)
		if got != c.want {
			t.Errorf("FindDuplicate(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}
