package main

import "testing"

func TestMaxSubArray_Classic(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	if got := MaxSubArray(nums); got != 6 {
		t.Errorf("MaxSubArray(%v) = %d, want 6", nums, got)
	}
}

func TestMaxSubArray_AllNegative(t *testing.T) {
	nums := []int{-3, -2, -1}
	if got := MaxSubArray(nums); got != -1 {
		t.Errorf("MaxSubArray(%v) = %d, want -1", nums, got)
	}
}

func TestMaxSubArray_AllPositive(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	if got := MaxSubArray(nums); got != 10 {
		t.Errorf("MaxSubArray(%v) = %d, want 10", nums, got)
	}
}

func TestMaxSubArray_SingleElement(t *testing.T) {
	nums := []int{5}
	if got := MaxSubArray(nums); got != 5 {
		t.Errorf("MaxSubArray(%v) = %d, want 5", nums, got)
	}
}

func TestMaxSubArray_LargeNegativeInMiddle(t *testing.T) {
	nums := []int{5, 4, -20, 7, 8}
	if got := MaxSubArray(nums); got != 15 {
		t.Errorf("MaxSubArray(%v) = %d, want 15", nums, got)
	}
}
