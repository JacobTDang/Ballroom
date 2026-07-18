package main

// SettlesBill returns whether amounts sums to bill. Currently
// unreliable -- find and fix the bug.
func SettlesBill(amounts []float64, bill float64) bool {
	total := 0.0
	for _, amount := range amounts {
		total += amount
	}
	return total == bill
}
