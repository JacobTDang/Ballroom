package main

import "testing"

func TestCarFleet(t *testing.T) {
	cases := []struct {
		target   int
		position []int
		speed    []int
		want     int
	}{
		{12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}, 3},
		{10, []int{3}, []int{3}, 1},
		{100, []int{0, 2, 4}, []int{4, 2, 1}, 1},
		{10, []int{0, 4, 8}, []int{1, 1, 1}, 3},
		{10, []int{0, 3, 6}, []int{5, 5, 5}, 3},
		{20, []int{1, 4}, []int{2, 1}, 1},
		{5, []int{5}, []int{1}, 1},
	}

	for _, c := range cases {
		got := CarFleet(c.target, c.position, c.speed)
		if got != c.want {
			t.Errorf("CarFleet(%d, %v, %v) = %d, want %d", c.target, c.position, c.speed, got, c.want)
		}
	}
}
