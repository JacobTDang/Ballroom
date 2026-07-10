package main

import (
	"reflect"
	"sort"
	"testing"
)

// normalizeExact sorts the outer list of board configurations
// lexicographically -- each board's row order is fixed (can't
// reorder rows), only which boards appear (and in what outer order)
// is unconstrained.
func normalizeExact(boards [][]string) [][]string {
	out := make([][]string, len(boards))
	for i, b := range boards {
		out[i] = append([]string(nil), b...)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
	return out
}

func TestSolveNQueens(t *testing.T) {
	cases := []struct {
		n    int
		want [][]string
	}{
		{
			4,
			[][]string{
				{".Q..", "...Q", "Q...", "..Q."},
				{"..Q.", "Q...", "...Q", ".Q.."},
			},
		},
		{1, [][]string{{"Q"}}},
	}

	for _, c := range cases {
		got := normalizeExact(SolveNQueens(c.n))
		want := normalizeExact(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SolveNQueens(%d) = %v, want %v (order-independent)", c.n, got, want)
		}
	}
}

func TestSolveNQueens_NoSolutionsForTwoOrThree(t *testing.T) {
	if got := SolveNQueens(2); len(got) != 0 {
		t.Errorf("SolveNQueens(2) = %v, want empty (no solution exists)", got)
	}
	if got := SolveNQueens(3); len(got) != 0 {
		t.Errorf("SolveNQueens(3) = %v, want empty (no solution exists)", got)
	}
}
