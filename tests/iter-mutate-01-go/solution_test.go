package main

import (
	"reflect"
	"testing"
)

func TestRemoveValue(t *testing.T) {
	cases := []struct {
		in     []int
		target int
		want   []int
	}{
		{[]int{1, 2, 2, 2, 3}, 2, []int{1, 3}},
		{[]int{2, 2, 5, 2}, 2, []int{5}},
		{[]int{1, 3, 5}, 9, []int{1, 3, 5}},
		{[]int{2, 2, 2, 2}, 2, []int{}},
		{[]int{5, 2, 2, 5}, 2, []int{5, 5}},
		{[]int{2}, 2, []int{}},
		{[]int{}, 2, []int{}},
	}
	for _, c := range cases {
		got := RemoveValue(c.in, c.target)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("RemoveValue(%v, %d) = %v, want %v", c.in, c.target, got, c.want)
		}
	}
}
