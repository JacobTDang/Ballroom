package main

import "testing"

func TestJump_Classic(t *testing.T) {
	nums := []int{2, 3, 1, 1, 4}
	if got := Jump(nums); got != 2 {
		t.Errorf("Jump(%v) = %d, want 2", nums, got)
	}
}

func TestJump_SingleElement(t *testing.T) {
	nums := []int{0}
	if got := Jump(nums); got != 0 {
		t.Errorf("Jump(%v) = %d, want 0", nums, got)
	}
}

func TestJump_AllOnes(t *testing.T) {
	nums := []int{1, 1, 1, 1}
	if got := Jump(nums); got != 3 {
		t.Errorf("Jump(%v) = %d, want 3", nums, got)
	}
}

func TestJump_BigFirstJump(t *testing.T) {
	nums := []int{5, 0, 0, 0, 0}
	if got := Jump(nums); got != 1 {
		t.Errorf("Jump(%v) = %d, want 1", nums, got)
	}
}

func TestJump_TwoElements(t *testing.T) {
	nums := []int{2, 1}
	if got := Jump(nums); got != 1 {
		t.Errorf("Jump(%v) = %d, want 1", nums, got)
	}
}
