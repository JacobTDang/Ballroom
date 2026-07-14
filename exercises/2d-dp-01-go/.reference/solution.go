package main

// UniquePaths returns the number of unique paths from the top-left to
// the bottom-right of an m x n grid, moving only right or down.
func UniquePaths(m int, n int) int {
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			dp[c] += dp[c-1]
		}
	}
	return dp[n-1]
}
