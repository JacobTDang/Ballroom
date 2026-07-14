package main

// FindTargetSumWays counts the number of ways to assign + or - to
// each num in nums so that the resulting expression evaluates to
// target.
func FindTargetSumWays(nums []int, target int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	if target > total || target < -total {
		return 0
	}
	if (total+target)%2 != 0 {
		return 0
	}
	subsetSum := (total + target) / 2

	dp := make([]int, subsetSum+1)
	dp[0] = 1
	for _, n := range nums {
		for s := subsetSum; s >= n; s-- {
			dp[s] += dp[s-n]
		}
	}
	return dp[subsetSum]
}
