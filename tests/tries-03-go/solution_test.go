package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindWords(t *testing.T) {
	cases := []struct {
		board [][]byte
		words []string
		want  []string
	}{
		{
			[][]byte{
				[]byte("oaan"),
				[]byte("etae"),
				[]byte("ihkr"),
				[]byte("iflv"),
			},
			[]string{"oath", "pea", "eat", "rain"},
			[]string{"eat", "oath"},
		},
		{
			[][]byte{
				[]byte("ab"),
				[]byte("cd"),
			},
			[]string{"abcb"},
			[]string{},
		},
		{
			[][]byte{
				[]byte("a"),
			},
			[]string{"a"},
			[]string{"a"},
		},
		{
			[][]byte{
				[]byte("aa"),
			},
			[]string{"aaa"},
			[]string{},
		},
		{
			[][]byte{
				[]byte("ab"),
				[]byte("cd"),
			},
			[]string{"abdc"},
			[]string{"abdc"},
		},
	}

	for _, c := range cases {
		got := FindWords(c.board, c.words)
		sort.Strings(got)
		want := append([]string(nil), c.want...)
		sort.Strings(want)
		if !reflect.DeepEqual(normalizeNilEmpty(got), normalizeNilEmpty(want)) {
			t.Errorf("FindWords(%v) = %v, want %v (order-independent)", c.words, got, want)
		}
	}
}

// normalizeNilEmpty treats a nil slice and an empty slice as equal —
// FindWords may return either for "no matches".
func normalizeNilEmpty(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return s
}
