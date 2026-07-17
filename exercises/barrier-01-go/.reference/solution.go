package main

import "sync"

// Barrier: a generation counter is what makes reuse safe -- each round
// waits on its own generation, so a round-2 arrival can never consume
// a round-1 release. The last arriver flips the generation and wakes
// everyone; waiters loop until THEIR generation has passed.
type Barrier struct {
	mu         sync.Mutex
	cond       *sync.Cond
	n          int
	arrived    int
	generation int
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	gen := b.generation
	b.arrived++
	if b.arrived == b.n {
		b.arrived = 0
		b.generation++
		b.cond.Broadcast()
		return
	}
	for gen == b.generation {
		b.cond.Wait()
	}
}
