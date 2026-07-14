package main

// LengthOfLIS returns the length of the longest strictly increasing
// subsequence of nums.
func LengthOfLIS(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	best := 0

	for i := 0; i < n; i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}

	return best
}
