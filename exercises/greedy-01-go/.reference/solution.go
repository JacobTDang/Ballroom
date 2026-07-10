package main

// MaxSubArray returns the largest sum of any contiguous subarray of nums.
func MaxSubArray(nums []int) int {
	best := nums[0]
	cur := nums[0]
	for _, n := range nums[1:] {
		if cur < 0 {
			cur = n
		} else {
			cur += n
		}
		if cur > best {
			best = cur
		}
	}
	return best
}
