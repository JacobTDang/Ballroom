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
		{[]int{2, 7}, 9, []int{0, 1}},
		{[]int{-5, -3, -1}, -8, []int{0, 1}},
		{[]int{-1, 1}, 0, []int{0, 1}},
		{[]int{-1000000000, 1000000000}, 0, []int{0, 1}},
		{[]int{1000, 2000, 3000, 4000, 5000, 7, 6000, 7000, 20, 8000}, 27, []int{5, 8}},
		{[]int{1000000000, -999999999}, 1, []int{0, 1}},
	}
	for _, c := range cases {
		got := TwoSum(c.nums, c.target)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("TwoSum(%v, %d) = %v, want %v", c.nums, c.target, got, c.want)
		}
	}
}
