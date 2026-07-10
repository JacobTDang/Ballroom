package main

import "testing"

func TestMinCostConnectPoints_Classic(t *testing.T) {
	points := [][]int{{0, 0}, {2, 2}, {3, 10}, {5, 2}, {7, 0}}
	if got := MinCostConnectPoints(points); got != 20 {
		t.Errorf("MinCostConnectPoints(%v) = %d, want 20", points, got)
	}
}

func TestMinCostConnectPoints_ThreePoints(t *testing.T) {
	points := [][]int{{3, 12}, {-2, 5}, {-4, 1}}
	if got := MinCostConnectPoints(points); got != 18 {
		t.Errorf("MinCostConnectPoints(%v) = %d, want 18", points, got)
	}
}

func TestMinCostConnectPoints_SinglePoint(t *testing.T) {
	points := [][]int{{0, 0}}
	if got := MinCostConnectPoints(points); got != 0 {
		t.Errorf("MinCostConnectPoints(%v) = %d, want 0", points, got)
	}
}
