package main

// NumDecodings returns the number of ways to decode the digit string s
// using the 'A'-'Z' -> "1"-"26" mapping.
func NumDecodings(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}

	dp := make([]int, n+1)
	dp[0] = 1
	if s[0] != '0' {
		dp[1] = 1
	}

	for i := 2; i <= n; i++ {
		if s[i-1] != '0' {
			dp[i] += dp[i-1]
		}
		twoDigit := (int(s[i-2]-'0'))*10 + int(s[i-1]-'0')
		if twoDigit >= 10 && twoDigit <= 26 {
			dp[i] += dp[i-2]
		}
	}

	return dp[n]
}
