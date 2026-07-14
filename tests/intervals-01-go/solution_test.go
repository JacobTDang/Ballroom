package main

import (
	"reflect"
	"testing"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		name        string
		intervals   [][]int
		newInterval []int
		want        [][]int
	}{
		{
			name:        "merge into middle",
			intervals:   [][]int{{1, 3}, {6, 9}},
			newInterval: []int{2, 5},
			want:        [][]int{{1, 5}, {6, 9}},
		},
		{
			name:        "merge several",
			intervals:   [][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}},
			newInterval: []int{4, 8},
			want:        [][]int{{1, 2}, {3, 10}, {12, 16}},
		},
		{
			name:        "empty list",
			intervals:   [][]int{},
			newInterval: []int{5, 7},
			want:        [][]int{{5, 7}},
		},
		{
			name:        "insert after all, no overlap",
			intervals:   [][]int{{1, 5}},
			newInterval: []int{6, 8},
			want:        [][]int{{1, 5}, {6, 8}},
		},
		{
			name:        "insert before all, no overlap",
			intervals:   [][]int{{1, 5}},
			newInterval: []int{0, 0},
			want:        [][]int{{0, 0}, {1, 5}},
		},
		{
			name:        "new interval swallows everything",
			intervals:   [][]int{{2, 3}, {4, 5}, {6, 7}},
			newInterval: []int{1, 10},
			want:        [][]int{{1, 10}},
		},
		{
			name:        "touching intervals merge",
			intervals:   [][]int{{1, 2}},
			newInterval: []int{2, 3},
			want:        [][]int{{1, 3}},
		},
		{
			name:        "gap of one does not merge",
			intervals:   [][]int{{1, 2}},
			newInterval: []int{3, 4},
			want:        [][]int{{1, 2}, {3, 4}},
		},
		{
			name:        "identical interval",
			intervals:   [][]int{{3, 5}},
			newInterval: []int{3, 5},
			want:        [][]int{{3, 5}},
		},
		{
			name:        "new interval contained within existing",
			intervals:   [][]int{{1, 10}},
			newInterval: []int{3, 5},
			want:        [][]int{{1, 10}},
		},
		{
			name:        "larger list, merges several in the middle",
			intervals:   [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}, {13, 14}, {15, 16}, {17, 18}, {19, 20}},
			newInterval: []int{6, 15},
			want:        [][]int{{1, 2}, {3, 4}, {5, 16}, {17, 18}, {19, 20}},
		},
		{
			name:        "constraint boundary values",
			intervals:   [][]int{{50000, 60000}},
			newInterval: []int{0, 100000},
			want:        [][]int{{0, 100000}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Insert(c.intervals, c.newInterval)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Insert(%v, %v) = %v, want %v", c.intervals, c.newInterval, got, c.want)
			}
		})
	}
}
