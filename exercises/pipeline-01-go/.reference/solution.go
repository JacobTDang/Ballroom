package main

import "sync"

// FanOutIn: the WaitGroup knows when every worker is done, and only
// its watcher goroutine may close(out) -- close-by-the-last-writer is
// the whole trick. The collector then drains until close, on the
// calling goroutine, so returning can't race the appends.
func FanOutIn(inputs []int, workers int, stage func(int) int) []int {
	in := make(chan int)
	out := make(chan int)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				out <- stage(v)
			}
		}()
	}

	go func() {
		for _, v := range inputs {
			in <- v
		}
		close(in)
	}()
	go func() {
		wg.Wait()
		close(out)
	}()

	var results []int
	for v := range out {
		results = append(results, v)
	}
	return results
}
