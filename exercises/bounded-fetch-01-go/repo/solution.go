package main

import "sync"

// RunLimited runs every task, with at most `limit` executing
// concurrently.
//
// TODO: this version launches everything at once -- the limit is
// ignored entirely.
func RunLimited(tasks []func(), limit int) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}
	wg.Wait()
}
