package main

import "testing"

func TestReverse_Classic(t *testing.T) {
	if got := Reverse(123); got != 321 {
		t.Errorf("Reverse(123) = %d, want 321", got)
	}
}

func TestReverse_Negative(t *testing.T) {
	if got := Reverse(-123); got != -321 {
		t.Errorf("Reverse(-123) = %d, want -321", got)
	}
}

func TestReverse_TrailingZero(t *testing.T) {
	if got := Reverse(120); got != 21 {
		t.Errorf("Reverse(120) = %d, want 21", got)
	}
}

func TestReverse_OverflowPositive(t *testing.T) {
	if got := Reverse(1534236469); got != 0 {
		t.Errorf("Reverse(1534236469) = %d, want 0", got)
	}
}

func TestReverse_OverflowNegative(t *testing.T) {
	if got := Reverse(-2147483648); got != 0 {
		t.Errorf("Reverse(-2147483648) = %d, want 0", got)
	}
}

func TestReverse_Zero(t *testing.T) {
	if got := Reverse(0); got != 0 {
		t.Errorf("Reverse(0) = %d, want 0", got)
	}
}
