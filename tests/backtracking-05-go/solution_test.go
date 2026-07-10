package main

import (
	"reflect"
	"sort"
	"testing"
)

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

func TestCombinationSum2(t *testing.T) {
	cases := []struct {
		candidates []int
		target     int
		want       [][]int
	}{
		{
			[]int{10, 1, 2, 7, 6, 1, 5}, 8,
			[][]int{{1, 1, 6}, {1, 2, 5}, {1, 7}, {2, 6}},
		},
		{
			[]int{2, 5, 2, 1, 2}, 5,
			[][]int{{1, 2, 2}, {5}},
		},
	}

	for _, c := range cases {
		got := normalize(CombinationSum2(c.candidates, c.target))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("CombinationSum2(%v, %d) = %v, want %v (order-independent)", c.candidates, c.target, got, want)
		}
	}
}
