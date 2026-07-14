package main

// MaxProduct returns the largest product of any contiguous non-empty
// subarray of nums.
func MaxProduct(nums []int) int {
	result := nums[0]
	curMax, curMin := nums[0], nums[0]

	for _, n := range nums[1:] {
		if n < 0 {
			curMax, curMin = curMin, curMax
		}
		curMax = max(n, curMax*n)
		curMin = min(n, curMin*n)
		result = max(result, curMax)
	}

	return result
}
