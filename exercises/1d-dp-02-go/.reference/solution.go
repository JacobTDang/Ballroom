package main

// MinCostClimbingStairs returns the minimum cost to reach the top of a
// staircase where cost[i] is the cost of stepping on stair i, and you
// may start from step 0 or step 1 for free, climbing 1 or 2 steps at a
// time.
func MinCostClimbingStairs(cost []int) int {
	n := len(cost)
	prev, curr := 0, 0
	for i := 2; i <= n; i++ {
		next := curr + cost[i-1]
		if alt := prev + cost[i-2]; alt < next {
			next = alt
		}
		prev, curr = curr, next
	}
	return curr
}
