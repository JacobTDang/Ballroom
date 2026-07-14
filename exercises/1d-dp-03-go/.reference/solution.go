package main

// Rob returns the maximum amount of money that can be robbed from
// houses arranged in a line, given nums[i] is the money in house i,
// without robbing two adjacent houses.
func Rob(nums []int) int {
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
