package main

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	cases := []struct {
		nums   []int
		target int
		want   []int
	}{
		{[]int{2, 7, 11, 15}, 9, []int{0, 1}},
		{[]int{3, 2, 4}, 6, []int{1, 2}},
		{[]int{3, 3}, 6, []int{0, 1}},
		{[]int{1, 2, 3, 4, 5}, 9, []int{3, 4}},
		{[]int{-3, 4, 3, 90}, 0, []int{0, 2}},
		{[]int{0, 4, 3, 0}, 0, []int{0, 3}},
	}
	for _, c := range cases {
		got := TwoSum(c.nums, c.target)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("TwoSum(%v, %d) = %v, want %v", c.nums, c.target, got, c.want)
		}
	}
}
