package main

import "testing"

func TestMinCostClimbingStairs_Three(t *testing.T) {
	cost := []int{10, 15, 20}
	if got := MinCostClimbingStairs(cost); got != 15 {
		t.Errorf("MinCostClimbingStairs(%v) = %d, want 15", cost, got)
	}
}

func TestMinCostClimbingStairs_Ten(t *testing.T) {
	cost := []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}
	if got := MinCostClimbingStairs(cost); got != 6 {
		t.Errorf("MinCostClimbingStairs(%v) = %d, want 6", cost, got)
	}
}

func TestMinCostClimbingStairs_TwoEqual(t *testing.T) {
	cost := []int{0, 0}
	if got := MinCostClimbingStairs(cost); got != 0 {
		t.Errorf("MinCostClimbingStairs(%v) = %d, want 0", cost, got)
	}
}

func TestMinCostClimbingStairs_BoundaryMaxValues(t *testing.T) {
	cost := []int{999, 999}
	if got := MinCostClimbingStairs(cost); got != 999 {
		t.Errorf("MinCostClimbingStairs(%v) = %d, want 999", cost, got)
	}
}

func TestMinCostClimbingStairs_LargerAscending(t *testing.T) {
	cost := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if got := MinCostClimbingStairs(cost); got != 25 {
		t.Errorf("MinCostClimbingStairs(%v) = %d, want 25", cost, got)
	}
}
