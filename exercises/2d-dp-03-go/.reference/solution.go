package main

// MaxProfit returns the maximum profit achievable buying and selling
// prices with unlimited transactions, subject to a mandatory one-day
// cooldown after selling before buying again.
func MaxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	hold := -prices[0]
	sold := 0
	rest := 0
	for i := 1; i < len(prices); i++ {
		prevHold, prevSold, prevRest := hold, sold, rest
		hold = max(prevHold, prevRest-prices[i])
		sold = prevHold + prices[i]
		rest = max(prevRest, prevSold)
	}
	return max(sold, rest)
}
