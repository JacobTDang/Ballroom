package main

import "testing"

func TestLeastInterval(t *testing.T) {
	cases := []struct {
		tasks []byte
		n     int
		want  int
	}{
		{[]byte("AAABBB"), 2, 8},
		{[]byte("AAABBB"), 0, 6},
		{[]byte("AAAAAABCDEFG"), 2, 16},
		{[]byte("A"), 5, 1},
		{[]byte("AAAB"), 3, 9},
		{[]byte("AB"), 2, 2},
	}

	for _, c := range cases {
		got := LeastInterval(c.tasks, c.n)
		if got != c.want {
			t.Errorf("LeastInterval(%q, %d) = %d, want %d", c.tasks, c.n, got, c.want)
		}
	}
}
