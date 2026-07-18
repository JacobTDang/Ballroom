package main

import "testing"

func TestAlign(t *testing.T) {
	cases := []struct {
		in, k, want int
	}{
		{7, 4, 4},
		{8, 4, 8},
		{0, 4, 0},
		{-1, 4, -4},
		{-7, 4, -8}, // pinned: truncation gives -4, floor gives -8
		{-8, 4, -8}, // pinned: exact multiple -- must not overcorrect to -12
		{-9, 4, -12},
		{-10, 3, -12},
		{-12, 3, -12},
	}
	for _, c := range cases {
		got := Align(c.in, c.k)
		if got != c.want {
			t.Errorf("Align(%d, %d) = %d, want %d", c.in, c.k, got, c.want)
		}
	}
}
