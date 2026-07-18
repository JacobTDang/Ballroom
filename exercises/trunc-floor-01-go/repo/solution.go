package main

// Align rounds t down to the start of its k-wide bucket. Currently
// wrong for negative t -- find and fix the bug.
func Align(t, k int) int {
	return t / k * k
}
