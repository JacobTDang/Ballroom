package main

// MaxOf returns the largest value in v.
func MaxOf(v []int) int {
	best := v[0]
	for i := 1; i < len(v); i++ {
		if v[i] > best {
			best = v[i]
		}
	}
	return best
}
