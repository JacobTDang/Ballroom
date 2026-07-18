package main

import "strings"

// BuildLog joins log chunks that arrive newest-first into a single
// oldest-first log (no separator between chunks -- each chunk already
// carries its own formatting).
func BuildLog(chunks []string) string {
	var b strings.Builder
	for i := len(chunks) - 1; i >= 0; i-- {
		b.WriteString(chunks[i])
	}
	return b.String()
}
