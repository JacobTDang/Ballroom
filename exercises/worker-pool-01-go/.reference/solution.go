package main

import "sync"

// ProcessAll: a channel of job indices feeds `workers` goroutines;
// each writes results[i] for the i it took -- distinct indices, so the
// slice needs no lock. Order falls out for free because position, not
// completion time, decides where a result lands.
func ProcessAll(jobs []int, workers int, fn func(int) int) []int {
	results := make([]int, len(jobs))
	indices := make(chan int)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range indices {
				results[i] = fn(jobs[i])
			}
		}()
	}
	for i := range jobs {
		indices <- i
	}
	close(indices)
	wg.Wait()
	return results
}
