package main

import "testing"

func TestEvalRPN(t *testing.T) {
	cases := []struct {
		tokens []string
		want   int
	}{
		{[]string{"2", "1", "+", "3", "*"}, 9},
		{[]string{"4", "13", "5", "/", "+"}, 6},
		{[]string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}, 22},
		{[]string{"18"}, 18},
		{[]string{"4", "3", "-"}, 1},
		{[]string{"-3", "4", "+"}, 1},
		{[]string{"7", "-2", "/"}, -3},
		{[]string{"5", "5", "*", "5", "*"}, 125},
	}

	for _, c := range cases {
		got := EvalRPN(c.tokens)
		if got != c.want {
			t.Errorf("EvalRPN(%v) = %d, want %d", c.tokens, got, c.want)
		}
	}
}
