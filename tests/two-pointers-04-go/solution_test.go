package main

import "testing"

func TestMaxArea(t *testing.T) {
	cases := []struct {
		height []int
		want   int
	}{
		{[]int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{[]int{1, 1}, 1},
		{[]int{4, 3, 2, 1, 4}, 16},
		{[]int{1, 2, 1}, 2},
		{[]int{1, 2, 4, 3}, 4},
		{[]int{1, 3, 2, 5, 25, 24, 5}, 24},
	}

	for _, c := range cases {
		got := MaxArea(c.height)
		if got != c.want {
			t.Errorf("MaxArea(%v) = %d, want %d", c.height, got, c.want)
		}
	}
}
