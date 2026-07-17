package main

// FanOutIn fans inputs out to `workers` goroutines running stage and
// collects every result (order doesn't matter).
//
// TODO: this version returns without waiting for the collector to
// finish -- results go missing, and the append races with the return.
func FanOutIn(inputs []int, workers int, stage func(int) int) []int {
	in := make(chan int)
	out := make(chan int)

	for w := 0; w < workers; w++ {
		go func() {
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

	var results []int
	go func() {
		for i := 0; i < len(inputs); i++ {
			results = append(results, <-out)
		}
	}()
	return results
}
