package main

import "testing"

func TestRob_Classic(t *testing.T) {
	nums := []int{1, 2, 3, 1}
	if got := Rob(nums); got != 4 {
		t.Errorf("Rob(%v) = %d, want 4", nums, got)
	}
}

func TestRob_Larger(t *testing.T) {
	nums := []int{2, 7, 9, 3, 1}
	if got := Rob(nums); got != 12 {
		t.Errorf("Rob(%v) = %d, want 12", nums, got)
	}
}

func TestRob_SingleHouse(t *testing.T) {
	nums := []int{5}
	if got := Rob(nums); got != 5 {
		t.Errorf("Rob(%v) = %d, want 5", nums, got)
	}
}

func TestRob_TwoHouses(t *testing.T) {
	nums := []int{2, 1}
	if got := Rob(nums); got != 2 {
		t.Errorf("Rob(%v) = %d, want 2", nums, got)
	}
}

func TestRob_AllZeros(t *testing.T) {
	nums := []int{0, 0, 0, 0}
	if got := Rob(nums); got != 0 {
		t.Errorf("Rob(%v) = %d, want 0", nums, got)
	}
}

func TestRob_LargerMixedValues(t *testing.T) {
	nums := []int{5, 5, 10, 100, 10, 5}
	if got := Rob(nums); got != 110 {
		t.Errorf("Rob(%v) = %d, want 110", nums, got)
	}
}

func TestRob_BoundaryMaxValues(t *testing.T) {
	nums := []int{1000, 1000}
	if got := Rob(nums); got != 1000 {
		t.Errorf("Rob(%v) = %d, want 1000", nums, got)
	}
}
