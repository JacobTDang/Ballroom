package main

import "testing"

func TestMaxOf(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{3, 1, 4, 1, 5, 9, 2, 6}, 9},
		{[]int{-5, -1, -10}, -1},
		{[]int{42}, 42},
	}
	for _, c := range cases {
		got := MaxOf(c.in)
		if got != c.want {
			t.Errorf("MaxOf(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}
