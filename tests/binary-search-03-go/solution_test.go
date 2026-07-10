package main

import "testing"

func TestMinEatingSpeed(t *testing.T) {
	cases := []struct {
		piles []int
		h     int
		want  int
	}{
		{[]int{3, 6, 7, 11}, 8, 4},
		{[]int{30, 11, 23, 4, 20}, 5, 30},
		{[]int{30, 11, 23, 4, 20}, 6, 23},
		{[]int{1000000000}, 2, 500000000},
		{[]int{1}, 1, 1},
	}

	for _, c := range cases {
		got := MinEatingSpeed(c.piles, c.h)
		if got != c.want {
			t.Errorf("MinEatingSpeed(%v, %d) = %d, want %d", c.piles, c.h, got, c.want)
		}
	}
}
