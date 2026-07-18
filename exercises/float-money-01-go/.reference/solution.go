package main

import "math"

// SettlesBill returns whether amounts sums to bill, to the nearest
// cent.
func SettlesBill(amounts []float64, bill float64) bool {
	total := 0.0
	for _, amount := range amounts {
		total += amount
	}
	return math.Round(total*100) == math.Round(bill*100)
}
