package main

import "testing"

func TestEraseOverlapIntervals(t *testing.T) {
	cases := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{
			name:      "one overlap to remove",
			intervals: [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}},
			want:      1,
		},
		{
			name:      "all identical",
			intervals: [][]int{{1, 2}, {1, 2}, {1, 2}},
			want:      2,
		},
		{
			name:      "touching endpoints, no removal",
			intervals: [][]int{{1, 2}, {2, 3}},
			want:      0,
		},
		{
			name:      "single interval",
			intervals: [][]int{{1, 2}},
			want:      0,
		},
		{
			name:      "already non-overlapping",
			intervals: [][]int{{1, 2}, {3, 4}, {5, 6}},
			want:      0,
		},
		{
			name:      "heavy overlap needs two removals",
			intervals: [][]int{{1, 100}, {11, 22}, {1, 11}, {2, 12}},
			want:      2,
		},
		{
			name:      "boundary values, touching not overlapping",
			intervals: [][]int{{-50000, -49999}, {-49999, 50000}},
			want:      0,
		},
		{
			name:      "all sharing start, most must go",
			intervals: [][]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}},
			want:      4,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := EraseOverlapIntervals(c.intervals); got != c.want {
				t.Errorf("EraseOverlapIntervals(%v) = %d, want %d", c.intervals, got, c.want)
			}
		})
	}
}
