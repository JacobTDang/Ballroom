package main

import (
	"reflect"
	"testing"
)

// normalizeStrs treats a nil slice the same as an empty one -- whether
// Decode returns nil or []string{} for "no strings" isn't part of the
// contract, only the actual strings are.
func normalizeStrs(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func TestEncodeDecode(t *testing.T) {
	cases := [][]string{
		{"neet", "code", "love", "you"},
		{},
		{""},
		{"", "", ""},
		{"a#b", "c##d", "5#hello"},
		{"hello world", "foo,bar", "123"},
		{"4#abcd", "hello"},
		{"#####"},
		{"xyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxy"},
		{"123", "456", "0"},
		{"", "a", "", "b"},
	}
	for _, strs := range cases {
		encoded := Encode(strs)
		got := normalizeStrs(Decode(encoded))
		want := normalizeStrs(strs)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Decode(Encode(%v)) = %v, want %v (encoded was %q)", strs, got, want, encoded)
		}
	}
}
