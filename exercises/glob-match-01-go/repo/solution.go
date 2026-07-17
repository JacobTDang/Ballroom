package main

// Match reports whether s matches pattern: * (any run), ? (exactly
// one), [a-c] (one from a set/range). Whole-string. An unclosed [
// makes the pattern match nothing.
//
// TODO: this handles only a lone "*" and literal equality -- no
// per-position wildcards, no classes, no backtracking.
func Match(pattern, s string) bool {
	if pattern == "*" {
		return true
	}
	return pattern == s
}
