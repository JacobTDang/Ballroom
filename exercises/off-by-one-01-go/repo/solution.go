package main

// MaxOf returns the largest value in v. Currently panics — find and fix
// the bug.
func MaxOf(v []int) int {
	best := v[0]
	for i := 0; i <= len(v); i++ {
		if v[i] > best {
			best = v[i]
		}
	}
	return best
}
