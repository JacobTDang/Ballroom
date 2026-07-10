package main

// LeastInterval returns the minimum number of CPU intervals needed
// to run every task, with identical tasks separated by at least n
// intervals.
func LeastInterval(tasks []byte, n int) int {
	var freq [26]int
	for _, t := range tasks {
		freq[t-'A']++
	}
	maxFreq := 0
	for _, f := range freq {
		if f > maxFreq {
			maxFreq = f
		}
	}
	maxCount := 0
	for _, f := range freq {
		if f == maxFreq {
			maxCount++
		}
	}

	// The most frequent task needs (maxFreq-1) gaps of size (n+1)
	// after it, plus however many other tasks are tied for most
	// frequent (they fill the last slot alongside it). If there are
	// enough other distinct tasks to fill every gap, the answer is
	// just len(tasks) -- no idling needed.
	frameSize := (maxFreq-1)*(n+1) + maxCount
	if len(tasks) > frameSize {
		return len(tasks)
	}
	return frameSize
}
