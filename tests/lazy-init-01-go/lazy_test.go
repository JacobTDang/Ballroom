package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestInitRunsExactlyOnceUnderContention(t *testing.T) {
	var initCalls int64
	l := NewLazy(func() int {
		atomic.AddInt64(&initCalls, 1)
		time.Sleep(10 * time.Millisecond) // expensive on purpose
		return 42
	})

	const callers = 50
	results := make([]int, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = l.Get()
		}(i)
	}
	wg.Wait()

	if n := atomic.LoadInt64(&initCalls); n != 1 {
		t.Fatalf("init ran %d times under contention, want exactly once", n)
	}
	for i, v := range results {
		if v != 42 {
			t.Fatalf("caller %d got %d, want the init result 42", i, v)
		}
	}
}

func TestSequentialCallsStillOnce(t *testing.T) {
	var initCalls int64
	l := NewLazy(func() int {
		atomic.AddInt64(&initCalls, 1)
		return 7
	})
	for i := 0; i < 5; i++ {
		if v := l.Get(); v != 7 {
			t.Fatalf("Get = %d, want 7", v)
		}
	}
	if initCalls != 1 {
		t.Fatalf("init ran %d times across sequential Gets, want once", initCalls)
	}
}
