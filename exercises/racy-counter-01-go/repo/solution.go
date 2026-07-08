package main

import "sync"

// Counter increments a shared counter n times, once per goroutine, and
// returns the final count. It should always return exactly n.
func Counter(n int) int {
	count := 0
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++
		}()
	}
	wg.Wait()
	return count
}
