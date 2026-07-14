package main

// ClimbStairs returns the number of distinct ways to climb a staircase
// of n steps, taking 1 or 2 steps at a time.
func ClimbStairs(n int) int {
	if n <= 2 {
		return n
	}
	prev, curr := 1, 2
	for i := 3; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}
