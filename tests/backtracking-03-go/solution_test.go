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

// normalizeExact is like normalize but does NOT sort within each
// inner list -- permutations are order-sensitive internally (e.g.
// [1,2] and [2,1] are different results), only the outer list's
// order is unconstrained.
func normalizeExact(lists [][]int) [][]int {
	out := make([][]int, len(lists))
	for i, l := range lists {
		out[i] = append([]int(nil), l...)
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

func TestPermute(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{
			[]int{1, 2, 3},
			[][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
		{[]int{0, 1}, [][]int{{0, 1}, {1, 0}}},
		{[]int{1}, [][]int{{1}}},
	}

	for _, c := range cases {
		got := normalizeExact(Permute(c.nums))
		want := normalizeExact(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Permute(%v) = %v, want %v (order-independent)", c.nums, got, want)
		}
	}
}

func TestPermute_CountAndContentMatchUnordered(t *testing.T) {
	got := Permute([]int{1, 2, 3})
	if len(got) != 6 {
		t.Fatalf("Permute count = %d, want 6", len(got))
	}
	// every permutation, regardless of internal order, must contain
	// the same multiset {1,2,3}
	for _, p := range got {
		n := normalize([][]int{p})[0]
		if !reflect.DeepEqual(n, []int{1, 2, 3}) {
			t.Errorf("permutation %v does not contain exactly {1,2,3}", p)
		}
	}
}
