package main

// Align rounds t down to the start of its k-wide bucket.
func Align(t, k int) int {
	q := t / k
	if t%k != 0 && t < 0 {
		q--
	}
	return q * k
}
