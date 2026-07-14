package main

import "testing"

func TestFindTargetSumWays_Classic(t *testing.T) {
	nums := []int{1, 1, 1, 1, 1}
	if got := FindTargetSumWays(nums, 3); got != 5 {
		t.Errorf("FindTargetSumWays(%v, 3) = %d, want 5", nums, got)
	}
}

func TestFindTargetSumWays_Single(t *testing.T) {
	nums := []int{1}
	if got := FindTargetSumWays(nums, 1); got != 1 {
		t.Errorf("FindTargetSumWays(%v, 1) = %d, want 1", nums, got)
	}
}

func TestFindTargetSumWays_Unreachable(t *testing.T) {
	nums := []int{1, 2, 3}
	if got := FindTargetSumWays(nums, 100); got != 0 {
		t.Errorf("FindTargetSumWays(%v, 100) = %d, want 0", nums, got)
	}
}

func TestFindTargetSumWays_Zeros(t *testing.T) {
	nums := []int{0, 0, 0, 0, 0, 0, 0, 0, 1}
	if got := FindTargetSumWays(nums, 1); got != 256 {
		t.Errorf("FindTargetSumWays(%v, 1) = %d, want 256", nums, got)
	}
}

func TestFindTargetSumWays_ZeroTarget(t *testing.T) {
	nums := []int{1, 1}
	if got := FindTargetSumWays(nums, 0); got != 2 {
		t.Errorf("FindTargetSumWays(%v, 0) = %d, want 2", nums, got)
	}
}

func TestFindTargetSumWays_NegativeTarget(t *testing.T) {
	nums := []int{1, 1, 1, 1, 1}
	if got := FindTargetSumWays(nums, -3); got != 5 {
		t.Errorf("FindTargetSumWays(%v, -3) = %d, want 5", nums, got)
	}
}

func TestFindTargetSumWays_MixedZeroTarget(t *testing.T) {
	nums := []int{1, 2, 1}
	if got := FindTargetSumWays(nums, 0); got != 2 {
		t.Errorf("FindTargetSumWays(%v, 0) = %d, want 2", nums, got)
	}
}

func TestFindTargetSumWays_ParityImpossible(t *testing.T) {
	nums := []int{1, 2, 3}
	if got := FindTargetSumWays(nums, 1); got != 0 {
		t.Errorf("FindTargetSumWays(%v, 1) = %d, want 0", nums, got)
	}
}
