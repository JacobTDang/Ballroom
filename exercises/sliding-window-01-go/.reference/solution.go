package main

// MaxProfit returns the maximum profit from buying on one day and
// selling on a later day, or 0 if no profit is possible.
func MaxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	minPrice := prices[0]
	best := 0
	for _, p := range prices[1:] {
		if p-minPrice > best {
			best = p - minPrice
		}
		if p < minPrice {
			minPrice = p
		}
	}
	return best
}
