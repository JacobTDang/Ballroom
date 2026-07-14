package main

import (
	"reflect"
	"sort"
	"testing"
)

// normalize sorts the outer list of partitions lexicographically --
// each partition's internal order is already fixed by left-to-right
// splitting, only which partitions appear (and in what outer order)
// is unconstrained.
func normalize(lists [][]string) [][]string {
	out := make([][]string, len(lists))
	for i, l := range lists {
		out[i] = append([]string(nil), l...)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
	return out
}

func TestPartition(t *testing.T) {
	cases := []struct {
		s    string
		want [][]string
	}{
		{"aab", [][]string{{"a", "a", "b"}, {"aa", "b"}}},
		{"a", [][]string{{"a"}}},
		{"aba", [][]string{{"a", "b", "a"}, {"aba"}}},
		{"aa", [][]string{{"a", "a"}, {"aa"}}},
		{"abcba", [][]string{{"a", "b", "c", "b", "a"}, {"a", "bcb", "a"}, {"abcba"}}},
	}

	for _, c := range cases {
		got := normalize(Partition(c.s))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Partition(%q) = %v, want %v (order-independent)", c.s, got, want)
		}
	}
}
