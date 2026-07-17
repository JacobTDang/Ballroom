package main

// BoundedQueue: a buffered channel is exactly a bounded blocking FIFO
// -- Put blocks when the buffer is full, Get blocks when it's empty,
// and the runtime handles every wake-up.
type BoundedQueue struct {
	ch chan int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	return &BoundedQueue{ch: make(chan int, capacity)}
}

func (q *BoundedQueue) Put(v int) {
	q.ch <- v
}

func (q *BoundedQueue) Get() int {
	return <-q.ch
}
