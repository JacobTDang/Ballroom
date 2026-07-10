package main

import "math"

// Reverse returns x with its digits reversed, or 0 if the reversed
// value falls outside the signed 32-bit integer range.
func Reverse(x int) int {
	var result int64
	for x != 0 {
		digit := x % 10
		x /= 10
		result = result*10 + int64(digit)
		if result < math.MinInt32 || result > math.MaxInt32 {
			return 0
		}
	}
	return int(result)
}
