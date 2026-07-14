package main

// IsInterleave reports whether s3 can be formed by interleaving s1
// and s2 while preserving each string's relative character order.
func IsInterleave(s1 string, s2 string, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}

	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for i := 1; i <= m; i++ {
		dp[i][0] = dp[i-1][0] && s1[i-1] == s3[i-1]
	}
	for j := 1; j <= n; j++ {
		dp[0][j] = dp[0][j-1] && s2[j-1] == s3[j-1]
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			dp[i][j] = (dp[i-1][j] && s1[i-1] == s3[i+j-1]) ||
				(dp[i][j-1] && s2[j-1] == s3[i+j-1])
		}
	}
	return dp[m][n]
}
