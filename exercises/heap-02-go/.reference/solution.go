package main

import "container/heap"

// intMaxHeap is a container/heap max-heap of ints.
type intMaxHeap []int

func (h intMaxHeap) Len() int            { return len(h) }
func (h intMaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h intMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intMaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// LastStoneWeight repeatedly smashes the two heaviest stones
// together and returns the weight of whatever stone (if any) remains.
func LastStoneWeight(stones []int) int {
	h := &intMaxHeap{}
	*h = append(*h, stones...)
	heap.Init(h)
	for h.Len() > 1 {
		a := heap.Pop(h).(int)
		b := heap.Pop(h).(int)
		if a != b {
			heap.Push(h, a-b)
		}
	}
	if h.Len() == 0 {
		return 0
	}
	return (*h)[0]
}
