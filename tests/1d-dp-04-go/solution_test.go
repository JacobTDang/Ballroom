package main

import "testing"

func TestRobCircular_Classic(t *testing.T) {
	nums := []int{2, 3, 2}
	if got := RobCircular(nums); got != 3 {
		t.Errorf("RobCircular(%v) = %d, want 3", nums, got)
	}
}

func TestRobCircular_FourHouses(t *testing.T) {
	nums := []int{1, 2, 3, 1}
	if got := RobCircular(nums); got != 4 {
		t.Errorf("RobCircular(%v) = %d, want 4", nums, got)
	}
}

func TestRobCircular_ThreeInARow(t *testing.T) {
	nums := []int{1, 2, 3}
	if got := RobCircular(nums); got != 3 {
		t.Errorf("RobCircular(%v) = %d, want 3", nums, got)
	}
}

func TestRobCircular_SingleHouse(t *testing.T) {
	nums := []int{5}
	if got := RobCircular(nums); got != 5 {
		t.Errorf("RobCircular(%v) = %d, want 5", nums, got)
	}
}

func TestRobCircular_TwoHouses(t *testing.T) {
	nums := []int{5, 10}
	if got := RobCircular(nums); got != 10 {
		t.Errorf("RobCircular(%v) = %d, want 10", nums, got)
	}
}

func TestRobCircular_AllZeros(t *testing.T) {
	nums := []int{0, 0, 0, 0}
	if got := RobCircular(nums); got != 0 {
		t.Errorf("RobCircular(%v) = %d, want 0", nums, got)
	}
}

func TestRobCircular_LargerAlternating(t *testing.T) {
	nums := []int{2, 3, 2, 3, 2, 3, 2}
	if got := RobCircular(nums); got != 9 {
		t.Errorf("RobCircular(%v) = %d, want 9", nums, got)
	}
}

func TestRobCircular_BoundaryMaxValues(t *testing.T) {
	nums := []int{1000, 1000, 1000}
	if got := RobCircular(nums); got != 1000 {
		t.Errorf("RobCircular(%v) = %d, want 1000", nums, got)
	}
}
