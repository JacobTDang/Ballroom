package main

// CanJump returns whether the last index of nums is reachable, where
// nums[i] is the maximum jump length from index i.
func CanJump(nums []int) bool {
	farthest := 0
	for i, n := range nums {
		if i > farthest {
			return false
		}
		if i+n > farthest {
			farthest = i + n
		}
	}
	return true
}
