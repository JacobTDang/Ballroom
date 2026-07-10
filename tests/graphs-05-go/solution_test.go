package main

import "testing"

func TestOrangesRotting(t *testing.T) {
	cases := []struct {
		grid [][]int
		want int
	}{
		{[][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{[][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{[][]int{{0, 2}}, 0},
		{[][]int{{0}}, 0},
	}

	for _, c := range cases {
		got := OrangesRotting(c.grid)
		if got != c.want {
			t.Errorf("OrangesRotting(...) = %d, want %d", got, c.want)
		}
	}
}
