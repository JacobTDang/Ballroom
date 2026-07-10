package main

// HammingWeight returns the number of set bits ('1's) in the binary
// representation of n.
func HammingWeight(n uint32) int {
	count := 0
	for n != 0 {
		n &= n - 1
		count++
	}
	return count
}
