package main

import "sync"

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

// Transfer: deadlock needs a cycle, and a fixed global order on lock
// acquisition (by account ID) makes cycles impossible -- both crossed
// transfers now take the same lock first, so one simply waits its
// turn. No global serialization: transfers touching disjoint accounts
// still run fully in parallel.
func Transfer(from, to *Account, amount int) bool {
	first, second := from, to
	if second.ID < first.ID {
		first, second = second, first
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	if from.balance < amount {
		return false
	}
	from.balance -= amount
	to.balance += amount
	return true
}
