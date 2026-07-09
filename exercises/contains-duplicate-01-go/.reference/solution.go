package main

// ContainsDuplicate returns true if any value appears at least twice in nums.
func ContainsDuplicate(nums []int) bool {
	seen := make(map[int]bool, len(nums))
	for _, n := range nums {
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}
