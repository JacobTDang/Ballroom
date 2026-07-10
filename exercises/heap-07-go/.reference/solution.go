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

// MedianFinder tracks the running median of a stream of integers,
// using two heaps that split the stream in half: small holds the
// lower half (max-heap, so its top is the largest of the low half)
// and large holds the upper half (min-heap, so its top is the
// smallest of the high half). Kept balanced within 1 of each other
// after every insert, so the median is always at the top of one (or
// both) heaps.
type MedianFinder struct {
	small *intMaxHeap
	large *intMinHeap
}

func NewMedianFinder() *MedianFinder {
	small := &intMaxHeap{}
	large := &intMinHeap{}
	heap.Init(small)
	heap.Init(large)
	return &MedianFinder{small: small, large: large}
}

func (mf *MedianFinder) AddNum(num int) {
	heap.Push(mf.small, num)
	heap.Push(mf.large, heap.Pop(mf.small))
	if mf.large.Len() > mf.small.Len() {
		heap.Push(mf.small, heap.Pop(mf.large))
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if mf.small.Len() > mf.large.Len() {
		return float64((*mf.small)[0])
	}
	return float64((*mf.small)[0]+(*mf.large)[0]) / 2.0
}
