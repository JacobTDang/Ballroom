package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestStopDrainsEverythingAccepted(t *testing.T) {
	var handled int64
	s := NewServer(4, func(v int) {
		time.Sleep(2 * time.Millisecond)
		atomic.AddInt64(&handled, 1)
	})

	accepted := 0
	for i := 0; i < 200; i++ {
		if s.Submit(i) {
			accepted++
		}
	}

	done := make(chan struct{})
	go func() { s.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Stop never returned -- deadlock or workers never drained")
	}

	if got := atomic.LoadInt64(&handled); got != int64(accepted) {
		t.Fatalf("%d jobs handled after Stop returned, want every accepted job (%d)", got, accepted)
	}
}

func TestSubmitRefusedAfterStop(t *testing.T) {
	var handled int64
	s := NewServer(2, func(v int) { atomic.AddInt64(&handled, 1) })
	s.Submit(1)
	s.Stop()
	before := atomic.LoadInt64(&handled)

	if s.Submit(2) {
		t.Fatal("Submit accepted a job after Stop returned")
	}
	time.Sleep(20 * time.Millisecond)
	if got := atomic.LoadInt64(&handled); got != before {
		t.Fatalf("handled count grew from %d to %d after Stop", before, got)
	}
}
