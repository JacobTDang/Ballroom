package main

// KthLargest tracks the kth largest value seen so far in a stream of
// integers.
type KthLargest struct {
	k int
}

func NewKthLargest(k int, nums []int) *KthLargest {
	return &KthLargest{k: k}
}

func (kl *KthLargest) Add(val int) int {
	// TODO: implement
	return 0
}
