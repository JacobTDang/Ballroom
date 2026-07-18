package main

import "errors"

// MaxAdjacentDiff returns the largest absolute difference between two
// adjacent elements in v. Returns an error if v has fewer than two
// elements. Currently panics on some inputs — find and fix the bug.
func MaxAdjacentDiff(v []int) (int, error) {
	if len(v) == 0 {
		return 0, errors.New("max adjacent diff: need at least two values")
	}
	best := abs(v[1] - v[0])
	for i := 1; i < len(v)-1; i++ {
		d := abs(v[i+1] - v[i])
		if d > best {
			best = d
		}
	}
	return best, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
