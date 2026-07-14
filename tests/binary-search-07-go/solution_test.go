package main

import (
	"math"
	"testing"
)

func TestFindMedianSortedArrays(t *testing.T) {
	cases := []struct {
		nums1, nums2 []int
		want         float64
	}{
		{[]int{1, 3}, []int{2}, 2.0},
		{[]int{1, 2}, []int{3, 4}, 2.5},
		{[]int{}, []int{1}, 1.0},
		{[]int{2}, []int{}, 2.0},
		{[]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10}, 5.5},
		{[]int{1, 2, 3}, []int{4, 5, 6, 7, 8, 9}, 5.0},
		{[]int{1}, []int{2, 3, 4, 5}, 3.0},
		{[]int{100, 200}, []int{1, 2, 3}, 3.0},
	}

	for _, c := range cases {
		got := FindMedianSortedArrays(c.nums1, c.nums2)
		if math.Abs(got-c.want) > 1e-5 {
			t.Errorf("FindMedianSortedArrays(%v, %v) = %v, want %v", c.nums1, c.nums2, got, c.want)
		}
	}
}
