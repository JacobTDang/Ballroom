package main

import "sync"

// BoundedQueue is a fixed-capacity FIFO shared between producers and
// consumers. Put must block while full; Get must block while empty.
type BoundedQueue struct {
	mu       sync.Mutex
	items    []int
	capacity int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	return &BoundedQueue{capacity: capacity}
}

// Put adds an item. TODO: it currently ignores the capacity entirely
// and never blocks.
func (q *BoundedQueue) Put(v int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, v)
}

// Get removes and returns the oldest item. TODO: it currently returns
// 0 when empty instead of waiting for a Put.
func (q *BoundedQueue) Get() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return 0
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v
}
