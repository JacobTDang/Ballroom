package main

// Jump returns the minimum number of jumps needed to reach the last
// index of nums, where nums[i] is the maximum jump length from index i.
func Jump(nums []int) int {
	jumps := 0
	curEnd := 0
	farthest := 0
	for i := 0; i < len(nums)-1; i++ {
		if i+nums[i] > farthest {
			farthest = i + nums[i]
		}
		if i == curEnd {
			jumps++
			curEnd = farthest
		}
	}
	return jumps
}
