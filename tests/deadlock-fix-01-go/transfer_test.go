package main

import (
	"sync"
	"testing"
	"time"
)

func TestCrossedTransfersDoNotDeadlock(t *testing.T) {
	a := NewAccount(1, 10000)
	b := NewAccount(2, 10000)

	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					if i%2 == 0 {
						Transfer(a, b, 1)
					} else {
						Transfer(b, a, 1)
					}
				}
			}(i)
		}
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("crossed transfers deadlocked (each direction holding one lock, waiting on the other)")
	}

	if total := a.Balance() + b.Balance(); total != 20000 {
		t.Fatalf("total balance %d after transfers, want 20000 conserved", total)
	}
}

func TestInsufficientFundsMovesNothing(t *testing.T) {
	a := NewAccount(1, 5)
	b := NewAccount(2, 0)
	if Transfer(a, b, 10) {
		t.Fatal("Transfer succeeded with insufficient funds")
	}
	if a.Balance() != 5 || b.Balance() != 0 {
		t.Fatalf("balances %d/%d after failed transfer, want 5/0 untouched", a.Balance(), b.Balance())
	}
}
