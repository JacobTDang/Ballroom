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

	allOnes := [][]int{{1, 1}, {1, 1}}
	if got := MaxAreaOfIsland(allOnes); got != 4 {
		t.Errorf("MaxAreaOfIsland(allOnes) = %d, want 4", got)
	}

	scattered := [][]int{{1, 0, 1}, {0, 0, 0}, {1, 0, 1}}
	if got := MaxAreaOfIsland(scattered); got != 1 {
		t.Errorf("MaxAreaOfIsland(scattered) = %d, want 1", got)
	}
}
