package main

import "testing"

func TestLastStoneWeight(t *testing.T) {
	cases := []struct {
		stones []int
		want   int
	}{
		{[]int{2, 7, 4, 1, 8, 1}, 1},
		{[]int{1}, 1},
		{[]int{1, 1}, 0},
		{[]int{1, 3}, 2},
		{[]int{2, 2}, 0},
		{[]int{10, 4, 2, 10}, 2},
		{[]int{1, 1, 1, 1}, 0},
	}

	for _, c := range cases {
		got := LastStoneWeight(c.stones)
		if got != c.want {
			t.Errorf("LastStoneWeight(%v) = %d, want %d", c.stones, got, c.want)
		}
	}
}
