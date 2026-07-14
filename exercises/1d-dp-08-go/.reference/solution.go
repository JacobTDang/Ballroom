package main

// CoinChange returns the fewest number of coins from coins (unlimited
// supply of each denomination) needed to make up amount, or -1 if
// amount cannot be made up by any combination of the coins.
func CoinChange(coins []int, amount int) int {
	sentinel := amount + 1

	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = sentinel
	}

	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}

	if dp[amount] == sentinel {
		return -1
	}
	return dp[amount]
}
