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
	}

	for _, c := range cases {
		got := TwoSum(c.numbers, c.target)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("TwoSum(%v, %d) = %v, want %v", c.numbers, c.target, got, c.want)
		}
	}
}
