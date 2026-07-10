package main

import "testing"

func TestTrap(t *testing.T) {
	cases := []struct {
		height []int
		want   int
	}{
		{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{[]int{4, 2, 0, 3, 2, 5}, 9},
		{[]int{}, 0},
		{[]int{1, 2, 3, 4, 5}, 0},
		{[]int{5, 4, 3, 2, 1}, 0},
		{[]int{3, 0, 3}, 3},
		{[]int{2, 0, 2}, 2},
	}

	for _, c := range cases {
		got := Trap(c.height)
		if got != c.want {
			t.Errorf("Trap(%v) = %d, want %d", c.height, got, c.want)
		}
	}
}
