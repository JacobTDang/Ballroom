package main

import "testing"

func TestCanJump_Classic(t *testing.T) {
	nums := []int{2, 3, 1, 1, 4}
	if got := CanJump(nums); got != true {
		t.Errorf("CanJump(%v) = %v, want true", nums, got)
	}
}

func TestCanJump_ClassicFalse(t *testing.T) {
	nums := []int{3, 2, 1, 0, 4}
	if got := CanJump(nums); got != false {
		t.Errorf("CanJump(%v) = %v, want false", nums, got)
	}
}

func TestCanJump_SingleElement(t *testing.T) {
	nums := []int{0}
	if got := CanJump(nums); got != true {
		t.Errorf("CanJump(%v) = %v, want true", nums, got)
	}
}

func TestCanJump_ZeroAtStartBlocks(t *testing.T) {
	nums := []int{0, 1}
	if got := CanJump(nums); got != false {
		t.Errorf("CanJump(%v) = %v, want false", nums, got)
	}
}

func TestCanJump_BigFirstJumpCoversRest(t *testing.T) {
	nums := []int{5, 0, 0, 0, 0}
	if got := CanJump(nums); got != true {
		t.Errorf("CanJump(%v) = %v, want true", nums, got)
	}
}
