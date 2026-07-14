package main

import (
	"math"
	"testing"
)

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 1e-5
}

func TestMedianFinder(t *testing.T) {
	mf := NewMedianFinder()
	mf.AddNum(1)
	mf.AddNum(2)
	if got := mf.FindMedian(); !closeEnough(got, 1.5) {
		t.Errorf("FindMedian() = %v, want 1.5", got)
	}
	mf.AddNum(3)
	if got := mf.FindMedian(); !closeEnough(got, 2.0) {
		t.Errorf("FindMedian() = %v, want 2.0", got)
	}
}

func TestMedianFinder_SingleElement(t *testing.T) {
	mf := NewMedianFinder()
	mf.AddNum(42)
	if got := mf.FindMedian(); !closeEnough(got, 42.0) {
		t.Errorf("FindMedian() = %v, want 42.0", got)
	}
}

func TestMedianFinder_OutOfOrderInserts(t *testing.T) {
	mf := NewMedianFinder()
	for _, n := range []int{5, 1, 9, 3, 7} {
		mf.AddNum(n)
	}
	// sorted: 1,3,5,7,9 -> median 5
	if got := mf.FindMedian(); !closeEnough(got, 5.0) {
		t.Errorf("FindMedian() = %v, want 5.0", got)
	}
	mf.AddNum(10)
	// sorted: 1,3,5,7,9,10 -> median (5+7)/2=6.0
	if got := mf.FindMedian(); !closeEnough(got, 6.0) {
		t.Errorf("FindMedian() = %v, want 6.0", got)
	}
}

func TestMedianFinder_NegativeValues(t *testing.T) {
	mf := NewMedianFinder()
	for _, n := range []int{-5, -1, -3} {
		mf.AddNum(n)
	}
	// sorted: -5,-3,-1 -> median -3
	if got := mf.FindMedian(); !closeEnough(got, -3.0) {
		t.Errorf("FindMedian() = %v, want -3.0", got)
	}
	mf.AddNum(-2)
	// sorted: -5,-3,-2,-1 -> median (-3+-2)/2=-2.5
	if got := mf.FindMedian(); !closeEnough(got, -2.5) {
		t.Errorf("FindMedian() = %v, want -2.5", got)
	}
}
