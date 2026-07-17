package main

import "time"

// TokenBucket is a rate limiter shared by many goroutines. Allow takes
// a token if available; Refill (called by an external ticker) adds
// tokens, clamped at capacity.
//
// TODO: check-then-decrement below is two separate steps -- under
// contention this hands out more tokens than exist.
type TokenBucket struct {
	capacity int
	tokens   int
}

func NewTokenBucket(capacity int) *TokenBucket {
	return &TokenBucket{capacity: capacity, tokens: capacity}
}

func (b *TokenBucket) Allow() bool {
	if b.tokens > 0 {
		time.Sleep(time.Microsecond) // simulated bookkeeping -- widens the race window
		b.tokens--
		return true
	}
	return false
}

func (b *TokenBucket) Refill(n int) {
	b.tokens += n
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}
