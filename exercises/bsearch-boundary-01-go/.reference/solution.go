package main

// FirstAtLeast returns the index of the first element in v that is >=
// target, or len(v) if every element is smaller.
func FirstAtLeast(v []int, target int) int {
	lo, hi := 0, len(v)
	for lo < hi {
		mid := (lo + hi) / 2
		if v[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
