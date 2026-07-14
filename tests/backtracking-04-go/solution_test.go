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

func TestSubsetsWithDup(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{
			[]int{1, 2, 2},
			[][]int{{}, {1}, {1, 2}, {1, 2, 2}, {2}, {2, 2}},
		},
		{[]int{0}, [][]int{{}, {0}}},
		{
			[]int{4, 4, 4, 1, 4},
			[][]int{{}, {1}, {1, 4}, {1, 4, 4}, {1, 4, 4, 4}, {1, 4, 4, 4, 4}, {4}, {4, 4}, {4, 4, 4}, {4, 4, 4, 4}},
		},
	}

	for _, c := range cases {
		got := normalize(SubsetsWithDup(c.nums))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SubsetsWithDup(%v) = %v, want %v (order-independent)", c.nums, got, want)
		}
	}
}
