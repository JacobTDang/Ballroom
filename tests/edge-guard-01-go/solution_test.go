package main

import "testing"

func TestMaxAdjacentDiff(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{3, 1, 4, 1, 5, 9, 2, 6}, 7},
		{[]int{5, 5}, 0},
		{[]int{-5, -1, -10}, 9},
		{[]int{1, 100}, 99},
	}
	for _, c := range cases {
		got, err := MaxAdjacentDiff(c.in)
		if err != nil {
			t.Errorf("MaxAdjacentDiff(%v) unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("MaxAdjacentDiff(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestMaxAdjacentDiffEmptyErrors(t *testing.T) {
	_, err := MaxAdjacentDiff([]int{})
	if err == nil {
		t.Fatal("MaxAdjacentDiff([]) should return an error, got nil")
	}
}

func TestMaxAdjacentDiffSingleElementErrors(t *testing.T) {
	_, err := MaxAdjacentDiff([]int{42})
	if err == nil {
		t.Fatal("MaxAdjacentDiff([42]) should return an error (fewer than two values), got nil")
	}
}
