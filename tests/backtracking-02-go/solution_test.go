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

func TestCombinationSum(t *testing.T) {
	cases := []struct {
		candidates []int
		target     int
		want       [][]int
	}{
		{[]int{2, 3, 6, 7}, 7, [][]int{{2, 2, 3}, {7}}},
		{[]int{2, 3, 5}, 8, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}},
		{[]int{2}, 1, [][]int{}},
	}

	for _, c := range cases {
		got := normalize(CombinationSum(c.candidates, c.target))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("CombinationSum(%v, %d) = %v, want %v (order-independent)", c.candidates, c.target, got, want)
		}
	}
}
