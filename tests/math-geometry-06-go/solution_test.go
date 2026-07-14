package main

import (
	"math"
	"testing"
)

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func TestMyPow_PositiveExponent(t *testing.T) {
	got := MyPow(2.0, 10)
	want := 1024.0
	if !approxEqual(got, want) {
		t.Errorf("MyPow(2.0, 10) = %v, want %v", got, want)
	}
}

func TestMyPow_FractionalBase(t *testing.T) {
	got := MyPow(2.1, 3)
	want := 9.261
	if !approxEqual(got, want) {
		t.Errorf("MyPow(2.1, 3) = %v, want %v", got, want)
	}
}

func TestMyPow_NegativeExponent(t *testing.T) {
	got := MyPow(2.0, -2)
	want := 0.25
	if !approxEqual(got, want) {
		t.Errorf("MyPow(2.0, -2) = %v, want %v", got, want)
	}
}

func TestMyPow_ZeroExponent(t *testing.T) {
	got := MyPow(0.5, 0)
	want := 1.0
	if !approxEqual(got, want) {
		t.Errorf("MyPow(0.5, 0) = %v, want %v", got, want)
	}
}

func TestMyPow_NegativeBase(t *testing.T) {
	got := MyPow(-2.0, 3)
	want := -8.0
	if !approxEqual(got, want) {
		t.Errorf("MyPow(-2.0, 3) = %v, want %v", got, want)
	}
}

func TestMyPow_MinInt32Exponent(t *testing.T) {
	// x = 1 keeps the expected result exact regardless of exponent
	// magnitude, while still exercising the negate-the-most-negative-
	// exponent overflow edge case.
	got := MyPow(1.0, math.MinInt32)
	want := 1.0
	if !approxEqual(got, want) {
		t.Errorf("MyPow(1.0, MinInt32) = %v, want %v", got, want)
	}
}

func TestMyPow_NegativeBaseNegativeExponent(t *testing.T) {
	got := MyPow(-2.0, -2)
	want := 0.25
	if !approxEqual(got, want) {
		t.Errorf("MyPow(-2.0, -2) = %v, want %v", got, want)
	}
}

func TestMyPow_LargerPositiveExponent(t *testing.T) {
	got := MyPow(3.0, 5)
	want := 243.0
	if !approxEqual(got, want) {
		t.Errorf("MyPow(3.0, 5) = %v, want %v", got, want)
	}
}

func TestMyPow_FractionalSquared(t *testing.T) {
	got := MyPow(1.5, 2)
	want := 2.25
	if !approxEqual(got, want) {
		t.Errorf("MyPow(1.5, 2) = %v, want %v", got, want)
	}
}
