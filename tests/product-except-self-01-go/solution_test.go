package main

import (
	"reflect"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	cases := []struct {
		nums []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{[]int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
		{[]int{2, 3}, []int{3, 2}},
		{[]int{5, 0, 0, 4}, []int{0, 0, 0, 0}},
		{[]int{1, 1, 1, 1}, []int{1, 1, 1, 1}},
	}
	for _, c := range cases {
		got := ProductExceptSelf(c.nums)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ProductExceptSelf(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}
