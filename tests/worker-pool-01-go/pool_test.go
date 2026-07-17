package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestResultsInInputOrder(t *testing.T) {
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i
	}
	got := ProcessAll(jobs, 8, func(v int) int {
		time.Sleep(2 * time.Millisecond)
		return v * 2
	})
	if len(got) != len(jobs) {
		t.Fatalf("got %d results, want %d", len(got), len(jobs))
	}
	for i, v := range got {
		if v != i*2 {
			t.Fatalf("results[%d] = %d, want %d -- results must be in input order", i, v, i*2)
		}
	}
}

func TestActuallyRunsInParallelWithinTheBound(t *testing.T) {
	var inFlight, highWater int64
	fn := func(v int) int {
		cur := atomic.AddInt64(&inFlight, 1)
		for {
			hw := atomic.LoadInt64(&highWater)
			if cur <= hw || atomic.CompareAndSwapInt64(&highWater, hw, cur) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt64(&inFlight, -1)
		return v
	}

	jobs := make([]int, 64)
	ProcessAll(jobs, 8, fn)

	hw := atomic.LoadInt64(&highWater)
	if hw < 2 {
		t.Fatalf("high-water mark = %d concurrent jobs, want at least 2 -- the pool never ran anything in parallel", hw)
	}
	if hw > 8 {
		t.Fatalf("high-water mark = %d concurrent jobs, want at most the 8 workers requested", hw)
	}
}
