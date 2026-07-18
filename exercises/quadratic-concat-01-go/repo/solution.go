package main

// BuildLog joins log chunks that arrive newest-first into a single
// oldest-first log (no separator between chunks -- each chunk already
// carries its own formatting). Currently far slower than it should be
// on a large page -- find and fix the bug.
func BuildLog(chunks []string) string {
	s := ""
	for _, c := range chunks {
		s = c + s
	}
	return s
}
