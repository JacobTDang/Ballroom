package main

import "testing"

func TestSearchMatrix(t *testing.T) {
	m := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	cases := []struct {
		matrix [][]int
		target int
		want   bool
	}{
		{m, 3, true},
		{m, 13, false},
		{[][]int{{1}}, 1, true},
		{[][]int{{1, 3}}, 3, true},
		{m, 60, true},
		{m, 0, false},
	}

	for _, c := range cases {
		got := SearchMatrix(c.matrix, c.target)
		if got != c.want {
			t.Errorf("SearchMatrix(%v, %d) = %v, want %v", c.matrix, c.target, got, c.want)
		}
	}
}
