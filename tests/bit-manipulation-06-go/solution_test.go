package main

import "testing"

func TestGetSum_Classic(t *testing.T) {
	if got := GetSum(1, 1); got != 2 {
		t.Errorf("GetSum(1, 1) = %d, want 2", got)
	}
}

func TestGetSum_PositivePositive(t *testing.T) {
	if got := GetSum(2, 3); got != 5 {
		t.Errorf("GetSum(2, 3) = %d, want 5", got)
	}
}

func TestGetSum_NegativePositiveCancel(t *testing.T) {
	if got := GetSum(-1, 1); got != 0 {
		t.Errorf("GetSum(-1, 1) = %d, want 0", got)
	}
}

func TestGetSum_TwoNegatives(t *testing.T) {
	if got := GetSum(-5, -7); got != -12 {
		t.Errorf("GetSum(-5, -7) = %d, want -12", got)
	}
}

func TestGetSum_WithZero(t *testing.T) {
	if got := GetSum(0, 0); got != 0 {
		t.Errorf("GetSum(0, 0) = %d, want 0", got)
	}
}

func TestGetSum_MaxInt32Bounds(t *testing.T) {
	if got := GetSum(2147483647, 0); got != 2147483647 {
		t.Errorf("GetSum(2147483647, 0) = %d, want 2147483647", got)
	}
	if got := GetSum(-2147483648, 0); got != -2147483648 {
		t.Errorf("GetSum(-2147483648, 0) = %d, want -2147483648", got)
	}
}
