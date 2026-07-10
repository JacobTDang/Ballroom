package main

// MissingNumber returns the one number in [0, n] missing from nums,
// where n is len(nums).
func MissingNumber(nums []int) int {
	result := len(nums)
	for i, v := range nums {
		result ^= i ^ v
	}
	return result
}
