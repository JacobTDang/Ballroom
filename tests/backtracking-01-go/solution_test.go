package main

import (
	"reflect"
	"sort"
	"testing"
)

// normalize sorts each inner list ascending, then sorts the outer
// list lexicographically -- Subsets' order isn't uniquely defined, so
// tests compare as a set of sets rather than an exact sequence.
func normalize(lists [][]int) [][]int {
	out := make([][]int, len(lists))
	for i, l := range lists {
		c := append([]int(nil), l...)
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

func TestSubsets(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{
			[]int{1, 2, 3},
			[][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		},
		{[]int{0}, [][]int{{}, {0}}},
		{[]int{1, 2}, [][]int{{}, {1}, {2}, {1, 2}}},
		{[]int{-1, 1}, [][]int{{}, {-1}, {1}, {-1, 1}}},
	}

	for _, c := range cases {
		got := normalize(Subsets(c.nums))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Subsets(%v) = %v, want %v (order-independent)", c.nums, got, want)
		}
	}
}
