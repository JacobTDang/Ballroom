package main

import "testing"

func TestMissingNumber_Classic(t *testing.T) {
	nums := []int{3, 0, 1}
	if got := MissingNumber(nums); got != 2 {
		t.Errorf("MissingNumber(%v) = %d, want 2", nums, got)
	}
}

func TestMissingNumber_MissingAtEnd(t *testing.T) {
	nums := []int{0, 1}
	if got := MissingNumber(nums); got != 2 {
		t.Errorf("MissingNumber(%v) = %d, want 2", nums, got)
	}
}

func TestMissingNumber_MissingAtStart(t *testing.T) {
	nums := []int{1, 2, 3}
	if got := MissingNumber(nums); got != 0 {
		t.Errorf("MissingNumber(%v) = %d, want 0", nums, got)
	}
}

func TestMissingNumber_SingleElement(t *testing.T) {
	nums := []int{0}
	if got := MissingNumber(nums); got != 1 {
		t.Errorf("MissingNumber(%v) = %d, want 1", nums, got)
	}
}
