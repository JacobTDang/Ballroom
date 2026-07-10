package main

import (
	"strconv"
	"strings"
)

// Encode encodes a list of strings into a single string that Decode can
// reconstruct exactly, including strings that contain any characters
// (digits, delimiters, etc). Each string is prefixed with its length and
// a '#' delimiter, so the delimiter itself can safely appear inside a
// string without ambiguity.
func Encode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		b.WriteString(strconv.Itoa(len(s)))
		b.WriteByte('#')
		b.WriteString(s)
	}
	return b.String()
}

// Decode reverses Encode.
func Decode(s string) []string {
	var result []string
	i := 0
	for i < len(s) {
		j := i
		for s[j] != '#' {
			j++
		}
		length, _ := strconv.Atoi(s[i:j])
		start := j + 1
		result = append(result, s[start:start+length])
		i = start + length
	}
	return result
}
