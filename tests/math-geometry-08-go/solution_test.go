package main

import "testing"

func TestDetectSquares_Classic(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{3, 10})
	ds.Add([]int{11, 2})
	ds.Add([]int{3, 2})

	if got := ds.Count([]int{11, 10}); got != 1 {
		t.Errorf("Count([11,10]) = %d, want 1", got)
	}
	if got := ds.Count([]int{14, 8}); got != 0 {
		t.Errorf("Count([14,8]) = %d, want 0", got)
	}

	ds.Add([]int{11, 2})
	if got := ds.Count([]int{11, 10}); got != 2 {
		t.Errorf("Count([11,10]) after duplicate add = %d, want 2", got)
	}
}

func TestDetectSquares_EmptyState(t *testing.T) {
	ds := NewDetectSquares()
	if got := ds.Count([]int{0, 0}); got != 0 {
		t.Errorf("Count([0,0]) on empty structure = %d, want 0", got)
	}
}

func TestDetectSquares_SymmetricBothSides(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{0, 2})
	ds.Add([]int{2, 0})
	ds.Add([]int{2, 2})
	ds.Add([]int{-2, 0})
	ds.Add([]int{-2, 2})

	if got := ds.Count([]int{0, 0}); got != 2 {
		t.Errorf("Count([0,0]) = %d, want 2", got)
	}
}

func TestDetectSquares_OneSidedOnly(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{1, 4})
	ds.Add([]int{4, 1})
	ds.Add([]int{4, 4})

	if got := ds.Count([]int{1, 1}); got != 1 {
		t.Errorf("Count([1,1]) = %d, want 1", got)
	}
}

func TestDetectSquares_DuplicateFrequencyMultiplication(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{1, 4})
	ds.Add([]int{1, 4})
	ds.Add([]int{1, 4})
	ds.Add([]int{4, 1})
	ds.Add([]int{4, 1})
	ds.Add([]int{4, 4})

	if got := ds.Count([]int{1, 1}); got != 6 {
		t.Errorf("Count([1,1]) = %d, want 6", got)
	}
}

func TestDetectSquares_NoMatchingXCoordinate(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{5, 5})
	ds.Add([]int{5, 9})
	ds.Add([]int{9, 5})
	ds.Add([]int{9, 9})

	if got := ds.Count([]int{100, 100}); got != 0 {
		t.Errorf("Count([100,100]) = %d, want 0", got)
	}
}

func TestDetectSquares_CountDoesNotMutateState(t *testing.T) {
	ds := NewDetectSquares()
	ds.Add([]int{0, 2})
	ds.Add([]int{2, 0})
	ds.Add([]int{2, 2})
	ds.Add([]int{-2, 0})
	ds.Add([]int{-2, 2})

	if got := ds.Count([]int{0, 0}); got != 2 {
		t.Errorf("Count([0,0]) first call = %d, want 2", got)
	}
	if got := ds.Count([]int{0, 0}); got != 2 {
		t.Errorf("Count([0,0]) second call = %d, want 2", got)
	}

	ds.Add([]int{0, 2})
	if got := ds.Count([]int{0, 0}); got != 4 {
		t.Errorf("Count([0,0]) after re-add = %d, want 4", got)
	}
}
