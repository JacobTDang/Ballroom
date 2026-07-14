package main

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	cases := []struct {
		numbers []int
		target  int
		want    []int
	}{
		{[]int{2, 7, 11, 15}, 9, []int{1, 2}},
		{[]int{2, 3, 4}, 6, []int{1, 3}},
		{[]int{-1, 0}, -1, []int{1, 2}},
		{[]int{3, 3}, 6, []int{1, 2}},                     // minimum size, duplicate values
		{[]int{1, 2, 3, 4, 4, 9, 56, 90}, 8, []int{4, 5}}, // duplicates inside a larger array
		{[]int{-8, -3, 0, 4, 9, 13}, 5, []int{1, 6}},      // first and last elements
		{[]int{-8, -3, 0, 4, 9, 13}, 4, []int{3, 4}},      // adjacent middle elements
		{[]int{-3, -1, 0, 2, 4, 5}, 1, []int{1, 5}},       // negative and positive pairing
		{[]int{5, 25, 75}, 100, []int{2, 3}},              // last two elements
	}

	for _, c := range cases {
		got := TwoSum(c.numbers, c.target)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("TwoSum(%v, %d) = %v, want %v", c.numbers, c.target, got, c.want)
		}
	}
}
