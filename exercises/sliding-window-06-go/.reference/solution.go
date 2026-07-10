package main

// MaxSlidingWindow returns the maximum of each contiguous window of
// size k as it slides from the start of nums to the end.
func MaxSlidingWindow(nums []int, k int) []int {
	var deque []int // indices into nums, values strictly decreasing
	var res []int
	for i, n := range nums {
		for len(deque) > 0 && nums[deque[len(deque)-1]] < n {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		if deque[0] <= i-k {
			deque = deque[1:]
		}
		if i >= k-1 {
			res = append(res, nums[deque[0]])
		}
	}
	return res
}
