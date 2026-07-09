package main

import (
	"sync"
	"sync/atomic"
)

// Counter increments a shared counter n times, once per goroutine, and
// returns the final count. It always returns exactly n.
func Counter(n int) int {
	var count int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1)
		}()
	}
	wg.Wait()
	return int(count)
}
