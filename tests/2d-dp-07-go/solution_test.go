package main

import "testing"

func TestLongestIncreasingPath_Classic(t *testing.T) {
	matrix := [][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}
	if got := LongestIncreasingPath(matrix); got != 4 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 4", matrix, got)
	}
}

func TestLongestIncreasingPath_SecondClassic(t *testing.T) {
	matrix := [][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}
	if got := LongestIncreasingPath(matrix); got != 4 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 4", matrix, got)
	}
}

func TestLongestIncreasingPath_SingleCell(t *testing.T) {
	matrix := [][]int{{1}}
	if got := LongestIncreasingPath(matrix); got != 1 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 1", matrix, got)
	}
}

func TestLongestIncreasingPath_SingleRow(t *testing.T) {
	matrix := [][]int{{1, 2, 3, 4}}
	if got := LongestIncreasingPath(matrix); got != 4 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 4", matrix, got)
	}
}

func TestLongestIncreasingPath_AllEqual(t *testing.T) {
	matrix := [][]int{{7, 7}, {7, 7}}
	if got := LongestIncreasingPath(matrix); got != 1 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 1", matrix, got)
	}
}

func TestLongestIncreasingPath_SingleColumn(t *testing.T) {
	matrix := [][]int{{1}, {2}, {3}}
	if got := LongestIncreasingPath(matrix); got != 3 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 3", matrix, got)
	}
}

func TestLongestIncreasingPath_SnakeFullTraversal(t *testing.T) {
	matrix := [][]int{{1, 2, 3}, {6, 5, 4}, {7, 8, 9}}
	if got := LongestIncreasingPath(matrix); got != 9 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 9", matrix, got)
	}
}

func TestLongestIncreasingPath_NegativeValues(t *testing.T) {
	matrix := [][]int{{-1, -2}, {-3, -4}}
	if got := LongestIncreasingPath(matrix); got != 3 {
		t.Errorf("LongestIncreasingPath(%v) = %d, want 3", matrix, got)
	}
}
