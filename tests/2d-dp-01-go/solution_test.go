package main

import "testing"

func TestUniquePaths_Classic(t *testing.T) {
	if got := UniquePaths(3, 7); got != 28 {
		t.Errorf("UniquePaths(3, 7) = %d, want 28", got)
	}
}

func TestUniquePaths_Small(t *testing.T) {
	if got := UniquePaths(3, 2); got != 3 {
		t.Errorf("UniquePaths(3, 2) = %d, want 3", got)
	}
}

func TestUniquePaths_SingleCell(t *testing.T) {
	if got := UniquePaths(1, 1); got != 1 {
		t.Errorf("UniquePaths(1, 1) = %d, want 1", got)
	}
}

func TestUniquePaths_SingleRow(t *testing.T) {
	if got := UniquePaths(1, 5); got != 1 {
		t.Errorf("UniquePaths(1, 5) = %d, want 1", got)
	}
}

func TestUniquePaths_SingleColumn(t *testing.T) {
	if got := UniquePaths(5, 1); got != 1 {
		t.Errorf("UniquePaths(5, 1) = %d, want 1", got)
	}
}

func TestUniquePaths_Square(t *testing.T) {
	if got := UniquePaths(3, 3); got != 6 {
		t.Errorf("UniquePaths(3, 3) = %d, want 6", got)
	}
}

func TestUniquePaths_LargerSquare(t *testing.T) {
	if got := UniquePaths(10, 10); got != 48620 {
		t.Errorf("UniquePaths(10, 10) = %d, want 48620", got)
	}
}

func TestUniquePaths_BoundaryMaxWithMinOther(t *testing.T) {
	if got := UniquePaths(2, 100); got != 100 {
		t.Errorf("UniquePaths(2, 100) = %d, want 100", got)
	}
}
