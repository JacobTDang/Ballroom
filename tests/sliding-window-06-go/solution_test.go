package main

import (
	"reflect"
	"testing"
)

func TestMaxSlidingWindow(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{[]int{1}, 1, []int{1}},
		{[]int{1, -1}, 1, []int{1, -1}},
		{[]int{9, 11}, 2, []int{11}},
		{[]int{4, -2}, 2, []int{4}},
		{[]int{1, 3, 1, 2, 0, 5}, 3, []int{3, 3, 2, 5}},
		{[]int{7, 2, 4}, 2, []int{7, 4}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{5}},
		{[]int{-7, -8, 7, 5, 7, 1, 6, 0}, 4, []int{7, 7, 7, 7, 7}},
	}

	for _, c := range cases {
		got := MaxSlidingWindow(c.nums, c.k)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("MaxSlidingWindow(%v, %d) = %v, want %v", c.nums, c.k, got, c.want)
		}
	}
}
