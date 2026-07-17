package main

import "sync"

// TokenBucket: one mutex makes check-and-take a single atomic step --
// the whole fix. Refill clamps under the same lock so a refill racing
// an Allow can't overshoot capacity either.
type TokenBucket struct {
	mu       sync.Mutex
	capacity int
	tokens   int
}

func NewTokenBucket(capacity int) *TokenBucket {
	return &TokenBucket{capacity: capacity, tokens: capacity}
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func (b *TokenBucket) Refill(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens += n
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}
