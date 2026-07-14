package main

// WordBreak returns whether s can be segmented into a space-separated
// sequence of one or more words from wordDict.
func WordBreak(s string, wordDict []string) bool {
	dict := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		dict[w] = true
	}

	n := len(s)
	dp := make([]bool, n+1)
	dp[0] = true

	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && dict[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}

	return dp[n]
}
