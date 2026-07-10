package main

// SingleNumber returns the element of nums that appears exactly once,
// given every other element appears exactly twice.
func SingleNumber(nums []int) int {
	result := 0
	for _, n := range nums {
		result ^= n
	}
	return result
}
