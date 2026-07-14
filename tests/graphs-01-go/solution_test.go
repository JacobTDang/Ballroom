package main

import "testing"

func gridOf(rows []string) [][]byte {
	g := make([][]byte, len(rows))
	for i, r := range rows {
		g[i] = []byte(r)
	}
	return g
}

func TestNumIslands(t *testing.T) {
	cases := []struct {
		grid [][]byte
		want int
	}{
		{gridOf([]string{"11110", "11010", "11000", "00000"}), 1},
		{gridOf([]string{"11000", "11000", "00100", "00011"}), 3},
		{gridOf([]string{"0"}), 0},
		{gridOf([]string{"1"}), 1},
		{gridOf([]string{"000", "000"}), 0},
		{gridOf([]string{"11", "11"}), 1},
	}

	for _, c := range cases {
		got := NumIslands(c.grid)
		if got != c.want {
			t.Errorf("NumIslands(...) = %d, want %d", got, c.want)
		}
	}
}
