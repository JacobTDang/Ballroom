package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestGenerateParenthesis(t *testing.T) {
	cases := []struct {
		n    int
		want []string
	}{
		{1, []string{"()"}},
		{2, []string{"(())", "()()"}},
		{3, []string{"((()))", "(()())", "(())()", "()(())", "()()()"}},
	}

	for _, c := range cases {
		got := GenerateParenthesis(c.n)
		sort.Strings(got)
		want := append([]string(nil), c.want...)
		sort.Strings(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GenerateParenthesis(%d) = %v, want %v (order-independent)", c.n, got, want)
		}
	}
}
