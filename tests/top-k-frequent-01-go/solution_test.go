package main

import (
	"sort"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{1, 1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{[]int{1}, 1, []int{1}},
		{[]int{4, 1, -1, 2, -1, 2, 3}, 2, []int{-1, 2}},
		{[]int{5, 5, 5, 5, 3, 3, 1}, 1, []int{5}},
		{[]int{1, 2, 3}, 3, []int{1, 2, 3}},
		{[]int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}, 2, []int{1, 2}},
		{[]int{-5, -5, -3, -3, -3, -1}, 1, []int{-3}},
		{[]int{7, 7, 7}, 1, []int{7}},
		{[]int{-10000, -10000, 10000}, 1, []int{-10000}},
	}
	for _, c := range cases {
		got := TopKFrequent(c.nums, c.k)
		gotSorted := append([]int(nil), got...)
		sort.Ints(gotSorted)
		wantSorted := append([]int(nil), c.want...)
		sort.Ints(wantSorted)

		if len(gotSorted) != len(wantSorted) {
			t.Errorf("TopKFrequent(%v, %d) = %v, want set %v", c.nums, c.k, got, c.want)
			continue
		}
		for i := range gotSorted {
			if gotSorted[i] != wantSorted[i] {
				t.Errorf("TopKFrequent(%v, %d) = %v, want set %v", c.nums, c.k, got, c.want)
				break
			}
		}
	}
}
