package main

import "container/heap"

// intMinHeap is a standard container/heap min-heap of ints.
type intMinHeap []int

func (h intMinHeap) Len() int            { return len(h) }
func (h intMinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intMinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// KthLargest tracks the kth largest value seen so far in a stream of
// integers, using a min-heap capped at size k -- the heap's smallest
// element (the top) is always the kth largest overall.
type KthLargest struct {
	k    int
	heap *intMinHeap
}

func NewKthLargest(k int, nums []int) *KthLargest {
	h := &intMinHeap{}
	heap.Init(h)
	kl := &KthLargest{k: k, heap: h}
	for _, n := range nums {
		kl.Add(n)
	}
	return kl
}

func (kl *KthLargest) Add(val int) int {
	heap.Push(kl.heap, val)
	if kl.heap.Len() > kl.k {
		heap.Pop(kl.heap)
	}
	return (*kl.heap)[0]
}
