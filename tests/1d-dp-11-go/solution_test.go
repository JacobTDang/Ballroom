package main

import "testing"

func TestLengthOfLIS_Classic(t *testing.T) {
	nums := []int{10, 9, 2, 5, 3, 7, 101, 18}
	if got := LengthOfLIS(nums); got != 4 {
		t.Errorf("LengthOfLIS(%v) = %d, want 4", nums, got)
	}
}

func TestLengthOfLIS_RepeatedDip(t *testing.T) {
	nums := []int{0, 1, 0, 3, 2, 3}
	if got := LengthOfLIS(nums); got != 4 {
		t.Errorf("LengthOfLIS(%v) = %d, want 4", nums, got)
	}
}

func TestLengthOfLIS_AllEqual(t *testing.T) {
	nums := []int{7, 7, 7, 7}
	if got := LengthOfLIS(nums); got != 1 {
		t.Errorf("LengthOfLIS(%v) = %d, want 1", nums, got)
	}
}

func TestLengthOfLIS_SingleElement(t *testing.T) {
	nums := []int{5}
	if got := LengthOfLIS(nums); got != 1 {
		t.Errorf("LengthOfLIS(%v) = %d, want 1", nums, got)
	}
}

func TestLengthOfLIS_StrictlyDecreasing(t *testing.T) {
	nums := []int{5, 4, 3, 2, 1}
	if got := LengthOfLIS(nums); got != 1 {
		t.Errorf("LengthOfLIS(%v) = %d, want 1", nums, got)
	}
}

func TestLengthOfLIS_StrictlyIncreasing(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	if got := LengthOfLIS(nums); got != 5 {
		t.Errorf("LengthOfLIS(%v) = %d, want 5", nums, got)
	}
}

func TestLengthOfLIS_NegativeValues(t *testing.T) {
	nums := []int{-1, -2, 0, 1, -3, 5}
	if got := LengthOfLIS(nums); got != 4 {
		t.Errorf("LengthOfLIS(%v) = %d, want 4", nums, got)
	}
}

func TestLengthOfLIS_BoundaryValues(t *testing.T) {
	nums := []int{-10000, 10000}
	if got := LengthOfLIS(nums); got != 2 {
		t.Errorf("LengthOfLIS(%v) = %d, want 2", nums, got)
	}
}

func TestLengthOfLIS_DuplicatesBreakStreak(t *testing.T) {
	nums := []int{3, 3, 3, 4, 4, 5}
	if got := LengthOfLIS(nums); got != 3 {
		t.Errorf("LengthOfLIS(%v) = %d, want 3", nums, got)
	}
}
