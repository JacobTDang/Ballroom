package main

// MyPow computes x raised to the power n in O(log n) time.
func MyPow(x float64, n int) float64 {
	// Use a 64-bit exponent so negating the most-negative n never
	// overflows.
	exp := int64(n)
	if exp < 0 {
		x = 1 / x
		exp = -exp
	}

	result := 1.0
	base := x
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
}
