package main

import "testing"

func TestMaxProduct_Classic(t *testing.T) {
	nums := []int{2, 3, -2, 4}
	if got := MaxProduct(nums); got != 6 {
		t.Errorf("MaxProduct(%v) = %d, want 6", nums, got)
	}
}

func TestMaxProduct_ZeroSplits(t *testing.T) {
	nums := []int{-2, 0, -1}
	if got := MaxProduct(nums); got != 0 {
		t.Errorf("MaxProduct(%v) = %d, want 0", nums, got)
	}
}

func TestMaxProduct_TwoNegativesFlip(t *testing.T) {
	nums := []int{-2, 3, -4}
	if got := MaxProduct(nums); got != 24 {
		t.Errorf("MaxProduct(%v) = %d, want 24", nums, got)
	}
}

func TestMaxProduct_SingleNegative(t *testing.T) {
	nums := []int{-5}
	if got := MaxProduct(nums); got != -5 {
		t.Errorf("MaxProduct(%v) = %d, want -5", nums, got)
	}
}

func TestMaxProduct_AllPositive(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	if got := MaxProduct(nums); got != 24 {
		t.Errorf("MaxProduct(%v) = %d, want 24", nums, got)
	}
}

func TestMaxProduct_SinglePositive(t *testing.T) {
	nums := []int{7}
	if got := MaxProduct(nums); got != 7 {
		t.Errorf("MaxProduct(%v) = %d, want 7", nums, got)
	}
}

func TestMaxProduct_MultipleZeroSplitIslands(t *testing.T) {
	nums := []int{0, 2, 0, 3, 0}
	if got := MaxProduct(nums); got != 3 {
		t.Errorf("MaxProduct(%v) = %d, want 3", nums, got)
	}
}

func TestMaxProduct_WholeArrayEvenNegatives(t *testing.T) {
	nums := []int{2, -3, 4, -5}
	if got := MaxProduct(nums); got != 120 {
		t.Errorf("MaxProduct(%v) = %d, want 120", nums, got)
	}
}
