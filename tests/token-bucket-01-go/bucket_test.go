package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func hammer(b *TokenBucket, callers int) int64 {
	var allowed int64
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow() {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}
	wg.Wait()
	return allowed
}

func TestExactlyCapacityAllowedUnderContention(t *testing.T) {
	b := NewTokenBucket(100)
	if got := hammer(b, 300); got != 100 {
		t.Fatalf("%d of 300 concurrent Allow calls succeeded, want exactly the capacity 100", got)
	}
}

func TestRefillGrantsExactlyThatMany(t *testing.T) {
	b := NewTokenBucket(100)
	hammer(b, 300) // drain
	b.Refill(40)
	if got := hammer(b, 200); got != 40 {
		t.Fatalf("%d Allow calls succeeded after Refill(40), want exactly 40", got)
	}
}

func TestRefillClampsAtCapacity(t *testing.T) {
	b := NewTokenBucket(50)
	b.Refill(1000)
	if got := hammer(b, 200); got != 50 {
		t.Fatalf("%d Allow calls succeeded after an over-refill, want the capacity 50", got)
	}
}
