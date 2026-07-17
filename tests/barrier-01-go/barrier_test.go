package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNoOneProceedsEarlyAcrossRounds: each participant increments the
// round's arrival counter before Wait and asserts after Wait that all
// n arrived -- anyone slipping through early sees a short count.
func TestNoOneProceedsEarlyAcrossRounds(t *testing.T) {
	const n, rounds = 4, 5
	b := NewBarrier(n)
	arrivals := make([]int64, rounds)
	failures := make(chan string, n*rounds)

	var wg sync.WaitGroup
	for p := 0; p < n; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				// Staggered arrivals widen the window for early release.
				time.Sleep(time.Duration(p) * 3 * time.Millisecond)
				atomic.AddInt64(&arrivals[r], 1)
				b.Wait()
				if got := atomic.LoadInt64(&arrivals[r]); got != n {
					failures <- "participant proceeded past round with only some arrivals"
					return
				}
			}
		}(p)
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("barrier deadlocked (participants never all released)")
	}
	close(failures)
	for f := range failures {
		t.Fatal(f)
	}
}
