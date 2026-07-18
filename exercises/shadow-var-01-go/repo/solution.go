package main

import "math"

// MaxBelowLimit returns the largest value in v that is <= limit, or -1
// if no value qualifies. Currently always returns -1 — find and fix
// the bug.
func MaxBelowLimit(v []int, limit int) int {
	result := math.MinInt
	for _, x := range v {
		if x <= limit && x > result {
			result := x
			if result == limit {
				break
			}
		}
	}
	if result == math.MinInt {
		return -1
	}
	return result
}
