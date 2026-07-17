package main

import (
	"sync"
	"time"
)

// Account is a bank account with its own lock.
type Account struct {
	ID      int
	mu      sync.Mutex
	balance int
}

func NewAccount(id, balance int) *Account {
	return &Account{ID: id, balance: balance}
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Transfer moves amount from one account to another, locking both.
//
// TODO: locking from-then-to deadlocks the moment two transfers cross
// (A->B and B->A each hold one lock and wait for the other). Fix the
// ordering -- don't just wrap everything in one global lock.
func Transfer(from, to *Account, amount int) bool {
	from.mu.Lock()
	time.Sleep(time.Millisecond) // bookkeeping -- widens the inversion window
	to.mu.Lock()
	defer from.mu.Unlock()
	defer to.mu.Unlock()

	if from.balance < amount {
		return false
	}
	from.balance -= amount
	to.balance += amount
	return true
}
