package main

// Change returns the number of combinations of coins (unlimited
// supply of each denomination) that sum to amount.
func Change(amount int, coins []int) int {
	dp := make([]int, amount+1)
	dp[0] = 1
	for _, coin := range coins {
		for x := coin; x <= amount; x++ {
			dp[x] += dp[x-coin]
		}
	}
	return dp[amount]
}
