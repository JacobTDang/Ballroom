package main

import "testing"

func TestCoinChange_Classic(t *testing.T) {
	coins := []int{1, 2, 5}
	if got := CoinChange(coins, 11); got != 3 {
		t.Errorf("CoinChange(%v, 11) = %d, want 3", coins, got)
	}
}

func TestCoinChange_Impossible(t *testing.T) {
	coins := []int{2}
	if got := CoinChange(coins, 3); got != -1 {
		t.Errorf("CoinChange(%v, 3) = %d, want -1", coins, got)
	}
}

func TestCoinChange_ZeroAmount(t *testing.T) {
	coins := []int{1}
	if got := CoinChange(coins, 0); got != 0 {
		t.Errorf("CoinChange(%v, 0) = %d, want 0", coins, got)
	}
}

func TestCoinChange_SingleCoinExact(t *testing.T) {
	coins := []int{3, 7}
	if got := CoinChange(coins, 6); got != 2 {
		t.Errorf("CoinChange(%v, 6) = %d, want 2", coins, got)
	}
}

func TestCoinChange_LargeAmountOnlyOnes(t *testing.T) {
	coins := []int{1}
	if got := CoinChange(coins, 10000); got != 10000 {
		t.Errorf("CoinChange(%v, 10000) = %d, want 10000", coins, got)
	}
}

func TestCoinChange_UnreachableAmount(t *testing.T) {
	coins := []int{3, 5}
	if got := CoinChange(coins, 7); got != -1 {
		t.Errorf("CoinChange(%v, 7) = %d, want -1", coins, got)
	}
}

func TestCoinChange_MixedDenominations(t *testing.T) {
	coins := []int{1, 5, 10, 25}
	if got := CoinChange(coins, 63); got != 6 {
		t.Errorf("CoinChange(%v, 63) = %d, want 6", coins, got)
	}
}
