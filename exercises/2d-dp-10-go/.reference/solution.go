package main

// MaxCoins returns the maximum coins obtainable by bursting all
// balloons in nums, where bursting balloon i yields
// nums[left] * nums[i] * nums[right] using the current neighbors.
func MaxCoins(nums []int) int {
	n := len(nums)
	balloons := make([]int, n+2)
	balloons[0], balloons[n+1] = 1, 1
	for i, v := range nums {
		balloons[i+1] = v
	}

	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}

	for length := 2; length <= n+1; length++ {
		for l := 0; l+length <= n+1; l++ {
			r := l + length
			for k := l + 1; k < r; k++ {
				coins := dp[l][k] + dp[k][r] + balloons[l]*balloons[k]*balloons[r]
				if coins > dp[l][r] {
					dp[l][r] = coins
				}
			}
		}
	}
	return dp[0][n+1]
}
