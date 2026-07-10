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
	}

	for _, c := range cases {
		got := MaxSlidingWindow(c.nums, c.k)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("MaxSlidingWindow(%v, %d) = %v, want %v", c.nums, c.k, got, c.want)
		}
	}
}
