package main

// RobCircular returns the maximum amount of money that can be robbed
// from houses arranged in a circle (house 0 and house n-1 are
// adjacent), given nums[i] is the money in house i, without robbing
// two adjacent houses.
func RobCircular(nums []int) int {
	n := len(nums)
	if n == 1 {
		return nums[0]
	}
	return max(robLinear(nums[:n-1]), robLinear(nums[1:]))
}

// robLinear is the House Robber I logic for a non-circular line of
// houses.
func robLinear(nums []int) int {
	prev, curr := 0, 0
	for _, n := range nums {
		next := curr
		if alt := prev + n; alt > next {
			next = alt
		}
		prev, curr = curr, next
	}
	return curr
}
