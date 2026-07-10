package main

import (
	"reflect"
	"sort"
	"testing"
)

func normalize(lists [][]int) [][]int {
	out := make([][]int, len(lists))
	copy(out, lists)
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

func TestPacificAtlantic(t *testing.T) {
	heights := [][]int{
		{1, 2, 2, 3, 5},
		{3, 2, 3, 4, 4},
		{2, 4, 5, 3, 1},
		{6, 7, 1, 4, 5},
		{5, 1, 1, 2, 4},
	}
	want := [][]int{{0, 4}, {1, 3}, {1, 4}, {2, 2}, {3, 0}, {3, 1}, {4, 0}}

	got := normalize(PacificAtlantic(heights))
	wantNorm := normalize(want)
	if !reflect.DeepEqual(got, wantNorm) {
		t.Errorf("PacificAtlantic(...) = %v, want %v (order-independent)", got, wantNorm)
	}
}

func TestPacificAtlantic_SingleCellFlowsToBoth(t *testing.T) {
	heights := [][]int{{1}}
	got := normalize(PacificAtlantic(heights))
	want := [][]int{{0, 0}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PacificAtlantic(single) = %v, want %v", got, want)
	}
}
