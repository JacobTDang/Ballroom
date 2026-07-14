package main

import "testing"

func TestSwimInWater_Small(t *testing.T) {
	grid := [][]int{{0, 2}, {1, 3}}
	if got := SwimInWater(grid); got != 3 {
		t.Errorf("SwimInWater(%v) = %d, want 3", grid, got)
	}
}

func TestSwimInWater_Larger(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4},
		{24, 23, 22, 21, 5},
		{12, 13, 14, 15, 16},
		{11, 17, 18, 19, 20},
		{10, 9, 8, 7, 6},
	}
	if got := SwimInWater(grid); got != 16 {
		t.Errorf("SwimInWater(%v) = %d, want 16", grid, got)
	}
}

func TestSwimInWater_SingleCell(t *testing.T) {
	grid := [][]int{{0}}
	if got := SwimInWater(grid); got != 0 {
		t.Errorf("SwimInWater(%v) = %d, want 0", grid, got)
	}
}

func TestSwimInWater_SpiralBlocksDirectPath(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{7, 8, 3},
		{6, 5, 4},
	}
	if got := SwimInWater(grid); got != 4 {
		t.Errorf("SwimInWater(%v) = %d, want 4", grid, got)
	}
}
