package main

import "testing"

func TestCountPaths(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{"two by two open", [][]int{{0, 0}, {0, 0}}, 2},
		{"trivial single cell", [][]int{{0}}, 1},
		{"destination blocked", [][]int{{0, 0}, {0, 1}}, 0},
		// Anti-overfit: the destination cell itself is open, but every
		// route to it is blocked. A fix that just hardcodes the
		// destination base case to 1 without preserving the
		// blocked-cell check must still get 0 here.
		{"fully blocked", [][]int{{0, 1}, {1, 0}}, 0},
		{"with obstacle", [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, 2},
		{"single row", [][]int{{0, 0, 0}}, 1},
	}
	for _, c := range cases {
		got := CountPaths(c.grid)
		if got != c.want {
			t.Errorf("%s: CountPaths(%v) = %d, want %d", c.name, c.grid, got, c.want)
		}
	}
}
