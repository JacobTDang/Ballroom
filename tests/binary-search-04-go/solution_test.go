package main

import "testing"

func TestFindMin(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{3, 4, 5, 1, 2}, 1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0},
		{[]int{11, 13, 15, 17}, 11},
		{[]int{2, 1}, 1},
		{[]int{1}, 1},
		{[]int{1, 2, 3, 4, 5}, 1},
		{[]int{15, 18, 2, 3, 6, 12}, 2},
	}

	for _, c := range cases {
		got := FindMin(c.nums)
		if got != c.want {
			t.Errorf("FindMin(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}
