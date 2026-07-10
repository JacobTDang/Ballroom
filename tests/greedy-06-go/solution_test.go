package main

import "testing"

func TestMergeTriplets_Classic(t *testing.T) {
	triplets := [][]int{{2, 5, 3}, {1, 8, 4}, {1, 7, 5}}
	target := []int{2, 7, 5}
	if got := MergeTriplets(triplets, target); got != true {
		t.Errorf("MergeTriplets(%v, %v) = %v, want true", triplets, target, got)
	}
}

func TestMergeTriplets_ClassicFalse(t *testing.T) {
	triplets := [][]int{{3, 4, 5}, {4, 5, 6}}
	target := []int{3, 2, 5}
	if got := MergeTriplets(triplets, target); got != false {
		t.Errorf("MergeTriplets(%v, %v) = %v, want false", triplets, target, got)
	}
}

func TestMergeTriplets_SingleExact(t *testing.T) {
	triplets := [][]int{{5, 5, 5}}
	target := []int{5, 5, 5}
	if got := MergeTriplets(triplets, target); got != true {
		t.Errorf("MergeTriplets(%v, %v) = %v, want true", triplets, target, got)
	}
}

func TestMergeTriplets_AllPoisoned(t *testing.T) {
	triplets := [][]int{{10, 1, 1}, {1, 10, 1}, {1, 1, 10}}
	target := []int{5, 5, 5}
	if got := MergeTriplets(triplets, target); got != false {
		t.Errorf("MergeTriplets(%v, %v) = %v, want false", triplets, target, got)
	}
}

func TestMergeTriplets_PartialMatch(t *testing.T) {
	triplets := [][]int{{2, 1, 1}, {1, 2, 1}}
	target := []int{2, 2, 2}
	if got := MergeTriplets(triplets, target); got != false {
		t.Errorf("MergeTriplets(%v, %v) = %v, want false", triplets, target, got)
	}
}
