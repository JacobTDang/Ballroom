package main

import "testing"

func TestMaxAreaOfIsland(t *testing.T) {
	grid := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
	}
	if got := MaxAreaOfIsland(grid); got != 6 {
		t.Errorf("MaxAreaOfIsland(...) = %d, want 6", got)
	}

	empty := [][]int{{0, 0, 0, 0, 0, 0, 0, 0}}
	if got := MaxAreaOfIsland(empty); got != 0 {
		t.Errorf("MaxAreaOfIsland(empty) = %d, want 0", got)
	}

	single := [][]int{{1}}
	if got := MaxAreaOfIsland(single); got != 1 {
		t.Errorf("MaxAreaOfIsland(single) = %d, want 1", got)
	}
}
