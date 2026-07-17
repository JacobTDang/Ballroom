package main

import "sync"

// Barrier makes n goroutines rendezvous: each Wait blocks until all n
// have arrived, then all proceed -- and the barrier must be reusable
// for the next round.
//
// TODO: releasing by closing a channel works exactly once -- the
// second round finds it already closed and nobody waits at all.
type Barrier struct {
	mu      sync.Mutex
	n       int
	arrived int
	release chan struct{}
}

func NewBarrier(n int) *Barrier {
	return &Barrier{n: n, release: make(chan struct{})}
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	b.arrived++
	if b.arrived == b.n {
		b.arrived = 0
		close(b.release)
		b.mu.Unlock()
		return
	}
	release := b.release
	b.mu.Unlock()
	<-release
}
