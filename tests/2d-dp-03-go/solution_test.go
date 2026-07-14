package main

import "testing"

func TestMaxProfit_Classic(t *testing.T) {
	prices := []int{1, 2, 3, 0, 2}
	if got := MaxProfit(prices); got != 3 {
		t.Errorf("MaxProfit(%v) = %d, want 3", prices, got)
	}
}

func TestMaxProfit_SingleDay(t *testing.T) {
	prices := []int{1}
	if got := MaxProfit(prices); got != 0 {
		t.Errorf("MaxProfit(%v) = %d, want 0", prices, got)
	}
}

func TestMaxProfit_MonotonicIncreasing(t *testing.T) {
	prices := []int{1, 2, 4}
	if got := MaxProfit(prices); got != 3 {
		t.Errorf("MaxProfit(%v) = %d, want 3", prices, got)
	}
}

func TestMaxProfit_Empty(t *testing.T) {
	if got := MaxProfit(nil); got != 0 {
		t.Errorf("MaxProfit(nil) = %d, want 0", got)
	}
}

func TestMaxProfit_MonotonicDecreasing(t *testing.T) {
	prices := []int{5, 4, 3, 2, 1}
	if got := MaxProfit(prices); got != 0 {
		t.Errorf("MaxProfit(%v) = %d, want 0", prices, got)
	}
}

func TestMaxProfit_TwoDaysProfit(t *testing.T) {
	prices := []int{1, 2}
	if got := MaxProfit(prices); got != 1 {
		t.Errorf("MaxProfit(%v) = %d, want 1", prices, got)
	}
}

func TestMaxProfit_CooldownForcesWait(t *testing.T) {
	prices := []int{1, 4, 2, 7}
	if got := MaxProfit(prices); got != 6 {
		t.Errorf("MaxProfit(%v) = %d, want 6", prices, got)
	}
}

func TestMaxProfit_LargerMultiTrade(t *testing.T) {
	prices := []int{6, 1, 3, 2, 4, 7}
	if got := MaxProfit(prices); got != 6 {
		t.Errorf("MaxProfit(%v) = %d, want 6", prices, got)
	}
}

func TestMaxProfit_BoundaryValues(t *testing.T) {
	prices := []int{10000, 1}
	if got := MaxProfit(prices); got != 0 {
		t.Errorf("MaxProfit(%v) = %d, want 0", prices, got)
	}
}
