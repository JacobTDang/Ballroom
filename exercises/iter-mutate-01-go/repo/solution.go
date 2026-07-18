package main

// RemoveValue removes every occurrence of target from v, in place, and
// returns it. Currently leaves some matches behind — find and fix the
// bug.
func RemoveValue(v []int, target int) []int {
	i := 0
	for i < len(v) {
		if v[i] == target {
			v = append(v[:i], v[i+1:]...)
		}
		i++
	}
	return v
}
