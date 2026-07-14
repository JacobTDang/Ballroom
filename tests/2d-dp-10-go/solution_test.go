package main

import "testing"

func TestMaxCoins_Classic(t *testing.T) {
	nums := []int{3, 1, 5, 8}
	if got := MaxCoins(nums); got != 167 {
		t.Errorf("MaxCoins(%v) = %d, want 167", nums, got)
	}
}

func TestMaxCoins_TwoBalloons(t *testing.T) {
	nums := []int{1, 5}
	if got := MaxCoins(nums); got != 10 {
		t.Errorf("MaxCoins(%v) = %d, want 10", nums, got)
	}
}

func TestMaxCoins_SingleBalloon(t *testing.T) {
	nums := []int{7}
	if got := MaxCoins(nums); got != 7 {
		t.Errorf("MaxCoins(%v) = %d, want 7", nums, got)
	}
}

func TestMaxCoins_Ones(t *testing.T) {
	nums := []int{1, 1}
	if got := MaxCoins(nums); got != 2 {
		t.Errorf("MaxCoins(%v) = %d, want 2", nums, got)
	}
}

func TestMaxCoins_AllOnesLarger(t *testing.T) {
	nums := []int{1, 1, 1, 1}
	if got := MaxCoins(nums); got != 4 {
		t.Errorf("MaxCoins(%v) = %d, want 4", nums, got)
	}
}

func TestMaxCoins_ZeroValueBalloon(t *testing.T) {
	nums := []int{3, 0, 5}
	if got := MaxCoins(nums); got != 20 {
		t.Errorf("MaxCoins(%v) = %d, want 20", nums, got)
	}
}

func TestMaxCoins_LargerAscending(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	if got := MaxCoins(nums); got != 110 {
		t.Errorf("MaxCoins(%v) = %d, want 110", nums, got)
	}
}

func TestMaxCoins_Descending(t *testing.T) {
	nums := []int{5, 3, 1}
	if got := MaxCoins(nums); got != 25 {
		t.Errorf("MaxCoins(%v) = %d, want 25", nums, got)
	}
}
