package main

import "testing"

func TestChange_Classic(t *testing.T) {
	coins := []int{1, 2, 5}
	if got := Change(5, coins); got != 4 {
		t.Errorf("Change(5, %v) = %d, want 4", coins, got)
	}
}

func TestChange_NoWay(t *testing.T) {
	coins := []int{2}
	if got := Change(3, coins); got != 0 {
		t.Errorf("Change(3, %v) = %d, want 0", coins, got)
	}
}

func TestChange_ZeroAmount(t *testing.T) {
	coins := []int{1, 2, 3}
	if got := Change(0, coins); got != 1 {
		t.Errorf("Change(0, %v) = %d, want 1", coins, got)
	}
}

func TestChange_ExactSingleCoin(t *testing.T) {
	coins := []int{10}
	if got := Change(10, coins); got != 1 {
		t.Errorf("Change(10, %v) = %d, want 1", coins, got)
	}
}

func TestChange_LargerAmount(t *testing.T) {
	coins := []int{1, 2, 5}
	if got := Change(10, coins); got != 10 {
		t.Errorf("Change(10, %v) = %d, want 10", coins, got)
	}
}

func TestChange_SingleCoinNoDivide(t *testing.T) {
	coins := []int{3}
	if got := Change(7, coins); got != 0 {
		t.Errorf("Change(7, %v) = %d, want 0", coins, got)
	}
}

func TestChange_MoreDenominations(t *testing.T) {
	coins := []int{2, 5, 3, 6}
	if got := Change(10, coins); got != 5 {
		t.Errorf("Change(10, %v) = %d, want 5", coins, got)
	}
}

func TestChange_BoundaryAmountSingleCoin(t *testing.T) {
	coins := []int{1}
	if got := Change(500, coins); got != 1 {
		t.Errorf("Change(500, %v) = %d, want 1", coins, got)
	}
}
