package main

import (
	"container/heap"
	"sort"
)

// candidate is one interval currently covering the sweep position,
// ordered in the heap by size (smallest interval on top).
type candidate struct {
	size int
	end  int
}

type candidateHeap []candidate

func (h candidateHeap) Len() int            { return len(h) }
func (h candidateHeap) Less(i, j int) bool  { return h[i].size < h[j].size }
func (h candidateHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *candidateHeap) Push(x interface{}) { *h = append(*h, x.(candidate)) }
func (h *candidateHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// MinInterval returns, for each query, the size of the smallest
// interval that contains it (left <= query <= right), or -1 if no
// interval contains it. The result is in the same order as queries.
func MinInterval(intervals [][]int, queries []int) []int {
	sorted := make([][]int, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})

	type queryPos struct {
		val int
		idx int
	}
	order := make([]queryPos, len(queries))
	for i, q := range queries {
		order[i] = queryPos{val: q, idx: i}
	}
	sort.Slice(order, func(i, j int) bool {
		return order[i].val < order[j].val
	})

	result := make([]int, len(queries))
	h := &candidateHeap{}
	heap.Init(h)

	i := 0
	for _, q := range order {
		for i < len(sorted) && sorted[i][0] <= q.val {
			left, right := sorted[i][0], sorted[i][1]
			heap.Push(h, candidate{size: right - left + 1, end: right})
			i++
		}

		for h.Len() > 0 && (*h)[0].end < q.val {
			heap.Pop(h)
		}

		if h.Len() > 0 {
			result[q.idx] = (*h)[0].size
		} else {
			result[q.idx] = -1
		}
	}

	return result
}
