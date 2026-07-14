package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// normalizeGroups sorts each group's strings and then sorts the list of
// groups, so results can be compared regardless of ordering — Go's map
// iteration order is intentionally randomized, and any correct grouping
// is valid no matter which order it comes back in.
func normalizeGroups(groups [][]string) [][]string {
	out := make([][]string, len(groups))
	for i, g := range groups {
		gc := append([]string(nil), g...)
		sort.Strings(gc)
		out[i] = gc
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.Join(out[i], ",") < strings.Join(out[j], ",")
	})
	return out
}

func TestGroupAnagrams(t *testing.T) {
	cases := []struct {
		strs []string
		want [][]string
	}{
		{[]string{"eat", "tea", "tan", "ate", "nat", "bat"}, [][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}}},
		{[]string{""}, [][]string{{""}}},
		{[]string{"a"}, [][]string{{"a"}}},
		{[]string{"abc", "bca", "cab", "xyz"}, [][]string{{"abc", "bca", "cab"}, {"xyz"}}},
		{[]string{"cat", "dog", "bird"}, [][]string{{"cat"}, {"dog"}, {"bird"}}},
		{[]string{"abc", "bca", "cab", "acb"}, [][]string{{"abc", "bca", "cab", "acb"}}},
		{[]string{"", "", ""}, [][]string{{"", "", ""}}},
		{[]string{"bat", "tab", "cat", "act", "dog", "god", "xyz"}, [][]string{{"bat", "tab"}, {"cat", "act"}, {"dog", "god"}, {"xyz"}}},
		{[]string{"a", "b", "a", "c", "b"}, [][]string{{"a", "a"}, {"b", "b"}, {"c"}}},
	}
	for _, c := range cases {
		got := normalizeGroups(GroupAnagrams(c.strs))
		want := normalizeGroups(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GroupAnagrams(%v) = %v, want %v (order-independent)", c.strs, got, want)
		}
	}
}
