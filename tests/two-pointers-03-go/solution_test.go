package main

import (
	"reflect"
	"sort"
	"testing"
)

// normalize sorts each triplet ascending, then sorts the list of
// triplets — 3Sum's valid outputs aren't uniquely ordered, so tests
// compare as sets rather than asserting an exact sequence.
func normalize(triplets [][]int) [][]int {
	out := make([][]int, len(triplets))
	for i, tr := range triplets {
		c := append([]int(nil), tr...)
		sort.Ints(c)
		out[i] = c
	}
	sort.Slice(out, func(i, j int) bool {
		for k := 0; k < len(out[i]) && k < len(out[j]); k++ {
			if out[i][k] != out[j][k] {
				return out[i][k] < out[j][k]
			}
		}
		return len(out[i]) < len(out[j])
	})
	return out
}

func TestThreeSum(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{[]int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{[]int{0, 1, 1}, [][]int{}},
		{[]int{0, 0, 0}, [][]int{{0, 0, 0}}},
		{[]int{}, [][]int{}},
		{[]int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{[]int{-2, 0, 1, 1, 2}, [][]int{{-2, 0, 2}, {-2, 1, 1}}},
		{[]int{-3, -2, -1}, [][]int{}},
		{[]int{1, 2, 3}, [][]int{}},
		{[]int{1, -1}, [][]int{}},
		{[]int{3, -2, 1, 0, -1, -3, 2, -2, 0}, [][]int{{-3, 0, 3}, {-3, 1, 2}, {-2, -1, 3}, {-2, 0, 2}, {-1, 0, 1}}},
	}

	for _, c := range cases {
		got := normalize(ThreeSum(c.nums))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ThreeSum(%v) = %v, want %v (order-independent)", c.nums, got, want)
		}
	}
}
