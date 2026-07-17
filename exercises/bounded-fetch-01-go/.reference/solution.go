package main

import "sync"

// RunLimited: a buffered channel is the semaphore -- acquire by
// sending before the task body, release by receiving after it. The
// bound covers execution, not just goroutine launch, because the send
// happens inside the goroutine, immediately before running the task.
func RunLimited(tasks []func(), limit int) {
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			t()
		}(task)
	}
	wg.Wait()
}
