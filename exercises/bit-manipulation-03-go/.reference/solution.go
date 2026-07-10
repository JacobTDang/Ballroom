package main

// CountBits returns a slice ans of length n+1 where ans[i] is the number
// of set bits in the binary representation of i.
func CountBits(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i>>1] + (i & 1)
	}
	return ans
}
