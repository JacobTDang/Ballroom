package main

import "testing"

func TestSingleNumber_Classic(t *testing.T) {
	nums := []int{2, 2, 1}
	if got := SingleNumber(nums); got != 1 {
		t.Errorf("SingleNumber(%v) = %d, want 1", nums, got)
	}
}

func TestSingleNumber_LongerMix(t *testing.T) {
	nums := []int{4, 1, 2, 1, 2}
	if got := SingleNumber(nums); got != 4 {
		t.Errorf("SingleNumber(%v) = %d, want 4", nums, got)
	}
}

func TestSingleNumber_SingleElement(t *testing.T) {
	nums := []int{7}
	if got := SingleNumber(nums); got != 7 {
		t.Errorf("SingleNumber(%v) = %d, want 7", nums, got)
	}
}

func TestSingleNumber_NegativeNumbers(t *testing.T) {
	nums := []int{-1, -1, -2}
	if got := SingleNumber(nums); got != -2 {
		t.Errorf("SingleNumber(%v) = %d, want -2", nums, got)
	}
}

func TestSingleNumber_BoundaryValues(t *testing.T) {
	nums := []int{1000, 1000, -1000}
	if got := SingleNumber(nums); got != -1000 {
		t.Errorf("SingleNumber(%v) = %d, want -1000", nums, got)
	}
}

func TestSingleNumber_LargerMixedSet(t *testing.T) {
	nums := []int{5, 3, 5, 4, 3}
	if got := SingleNumber(nums); got != 4 {
		t.Errorf("SingleNumber(%v) = %d, want 4", nums, got)
	}
}
