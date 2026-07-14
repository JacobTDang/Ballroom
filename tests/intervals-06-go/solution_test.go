package main

import (
	"reflect"
	"testing"
)

func TestMinInterval(t *testing.T) {
	cases := []struct {
		name      string
		intervals [][]int
		queries   []int
		want      []int
	}{
		{
			name:      "classic",
			intervals: [][]int{{1, 4}, {2, 4}, {3, 6}, {4, 4}},
			queries:   []int{2, 3, 4, 5},
			want:      []int{3, 3, 1, 4},
		},
		{
			name:      "some queries uncovered",
			intervals: [][]int{{2, 3}, {2, 5}, {1, 8}, {20, 25}},
			queries:   []int{2, 19, 5, 22},
			want:      []int{2, -1, 4, 6},
		},
		{
			name:      "boundary endpoints and a miss",
			intervals: [][]int{{1, 10}},
			queries:   []int{1, 10, 11},
			want:      []int{10, 10, -1},
		},
		{
			name:      "single point interval",
			intervals: [][]int{{5, 5}},
			queries:   []int{5},
			want:      []int{1},
		},
		{
			name:      "boundary constraint values",
			intervals: [][]int{{1, 10000000}},
			queries:   []int{1, 10000000},
			want:      []int{10000000, 10000000},
		},
		{
			name:      "nested intervals of varying size, multiple queries",
			intervals: [][]int{{1, 100}, {10, 20}, {15, 16}, {50, 60}},
			queries:   []int{15, 55, 99, 5},
			want:      []int{2, 11, 100, 100},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := MinInterval(c.intervals, c.queries)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("MinInterval(%v, %v) = %v, want %v", c.intervals, c.queries, got, c.want)
			}
		})
	}
}
