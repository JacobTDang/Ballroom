package main

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		name      string
		intervals [][]int
		want      [][]int
	}{
		{
			name:      "overlapping pair",
			intervals: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			want:      [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "touching endpoints merge",
			intervals: [][]int{{1, 4}, {4, 5}},
			want:      [][]int{{1, 5}},
		},
		{
			name:      "unsorted input",
			intervals: [][]int{{15, 18}, {2, 6}, {1, 3}, {8, 10}},
			want:      [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "single interval",
			intervals: [][]int{{1, 4}},
			want:      [][]int{{1, 4}},
		},
		{
			name:      "one interval fully contains another",
			intervals: [][]int{{1, 10}, {2, 3}, {4, 5}},
			want:      [][]int{{1, 10}},
		},
		{
			name:      "no overlaps at all",
			intervals: [][]int{{1, 2}, {3, 4}, {5, 6}},
			want:      [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:      "boundary values, no overlap",
			intervals: [][]int{{0, 1}, {9999, 10000}},
			want:      [][]int{{0, 1}, {9999, 10000}},
		},
		{
			name:      "larger input, multiple merge chains",
			intervals: [][]int{{1, 3}, {2, 4}, {5, 7}, {6, 8}, {10, 12}, {15, 20}, {18, 25}, {30, 31}},
			want:      [][]int{{1, 4}, {5, 8}, {10, 12}, {15, 25}, {30, 31}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Merge(c.intervals)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Merge(%v) = %v, want %v", c.intervals, got, c.want)
			}
		})
	}
}
