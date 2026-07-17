package main

import (
	"sync"
	"testing"
	"time"
)

func TestEveryItemArrivesExactlyOnce(t *testing.T) {
	q := NewBoundedQueue(4)
	const producers, perProducer, consumers = 3, 200, 3
	total := producers * perProducer

	var wg sync.WaitGroup
	for p := 0; p < producers; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			for i := 0; i < perProducer; i++ {
				q.Put(p*perProducer + i)
			}
		}(p)
	}

	var mu sync.Mutex
	seen := make(map[int]int)
	var cg sync.WaitGroup
	for c := 0; c < consumers; c++ {
		cg.Add(1)
		go func() {
			defer cg.Done()
			for i := 0; i < total/consumers; i++ {
				v := q.Get()
				mu.Lock()
				seen[v]++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	cg.Wait()

	if len(seen) != total {
		t.Fatalf("saw %d distinct items, want %d (lost or fabricated items)", len(seen), total)
	}
	for v, n := range seen {
		if n != 1 {
			t.Fatalf("item %d consumed %d times, want exactly once", v, n)
		}
	}
}

func TestGetBlocksUntilPut(t *testing.T) {
	q := NewBoundedQueue(2)
	got := make(chan int, 1)
	go func() { got <- q.Get() }()

	select {
	case v := <-got:
		t.Fatalf("Get on an empty queue returned %d immediately, want it to block", v)
	case <-time.After(50 * time.Millisecond):
	}

	q.Put(7)
	select {
	case v := <-got:
		if v != 7 {
			t.Fatalf("Get = %d after Put(7), want 7", v)
		}
	case <-time.After(time.Second):
		t.Fatal("Get never woke up after a Put")
	}
}

func TestPutBlocksUntilGet(t *testing.T) {
	q := NewBoundedQueue(2)
	q.Put(1)
	q.Put(2)

	done := make(chan struct{})
	go func() {
		q.Put(3)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Put on a full queue returned immediately, want it to block at capacity")
	case <-time.After(50 * time.Millisecond):
	}

	if v := q.Get(); v != 1 {
		t.Fatalf("Get = %d, want FIFO order (1 first)", v)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("blocked Put never completed after a Get freed a slot")
	}
}
