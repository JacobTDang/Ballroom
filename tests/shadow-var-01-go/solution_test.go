package main

import "testing"

func TestMaxBelowLimit(t *testing.T) {
	cases := []struct {
		in    []int
		limit int
		want  int
	}{
		{[]int{3, 7, 2, 9, 5}, 7, 7},
		{[]int{3, 7, 2, 9, 5}, 6, 5},
		{[]int{10, 20, 30}, 5, -1},
		{[]int{-5, -1, -10}, -2, -5},
		{[]int{5}, 10, 5},
		{[]int{15}, 10, -1},
	}
	for _, c := range cases {
		got := MaxBelowLimit(c.in, c.limit)
		if got != c.want {
			t.Errorf("MaxBelowLimit(%v, %d) = %d, want %d", c.in, c.limit, got, c.want)
		}
	}
}
