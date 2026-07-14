package main

import (
	"reflect"
	"testing"
)

func TestDailyTemperatures(t *testing.T) {
	cases := []struct {
		temperatures []int
		want         []int
	}{
		{[]int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{[]int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{[]int{30, 60, 90}, []int{1, 1, 0}},
		{[]int{80, 79, 78}, []int{0, 0, 0}},
		{[]int{75}, []int{0}},
		{[]int{55, 55, 55, 60}, []int{3, 2, 1, 0}},
	}

	for _, c := range cases {
		got := DailyTemperatures(c.temperatures)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("DailyTemperatures(%v) = %v, want %v", c.temperatures, got, c.want)
		}
	}
}
