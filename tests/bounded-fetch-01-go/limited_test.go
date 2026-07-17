package main

import (
	"sync/atomic"
	"testing"
	"time"
)

// instrumented returns n tasks that track how many run at once.
func instrumented(n int, ran, highWater *int64) []func() {
	var inFlight int64
	tasks := make([]func(), n)
	for i := range tasks {
		tasks[i] = func() {
			cur := atomic.AddInt64(&inFlight, 1)
			for {
				hw := atomic.LoadInt64(highWater)
				if cur <= hw || atomic.CompareAndSwapInt64(highWater, hw, cur) {
					break
				}
			}
			time.Sleep(15 * time.Millisecond)
			atomic.AddInt64(&inFlight, -1)
			atomic.AddInt64(ran, 1)
		}
	}
	return tasks
}

func TestBoundHoldsAndEverythingRuns(t *testing.T) {
	var ran, highWater int64
	RunLimited(instrumented(32, &ran, &highWater), 4)

	if ran != 32 {
		t.Fatalf("%d tasks ran, want all 32", ran)
	}
	if hw := atomic.LoadInt64(&highWater); hw > 4 {
		t.Fatalf("high-water mark %d, want at most the limit 4", hw)
	}
	if hw := atomic.LoadInt64(&highWater); hw < 2 {
		t.Fatalf("high-water mark %d, want real parallelism (at least 2) under limit 4", hw)
	}
}

func TestLimitOneIsSerial(t *testing.T) {
	var ran, highWater int64
	RunLimited(instrumented(6, &ran, &highWater), 1)
	if ran != 6 {
		t.Fatalf("%d tasks ran, want all 6", ran)
	}
	if hw := atomic.LoadInt64(&highWater); hw != 1 {
		t.Fatalf("high-water mark %d with limit 1, want exactly serial execution", hw)
	}
}

func TestLimitLargerThanTasks(t *testing.T) {
	var ran, highWater int64
	RunLimited(instrumented(3, &ran, &highWater), 10)
	if ran != 3 {
		t.Fatalf("%d tasks ran, want all 3", ran)
	}
}
