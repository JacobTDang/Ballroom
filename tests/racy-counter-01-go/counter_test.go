package main

import "testing"

func TestCounter(t *testing.T) {
	const n = 1000
	got := Counter(n)
	if got != n {
		t.Fatalf("Counter(%d) = %d, want %d", n, got, n)
	}
}

func TestCounter_SingleGoroutine(t *testing.T) {
	const n = 1
	got := Counter(n)
	if got != n {
		t.Fatalf("Counter(%d) = %d, want %d", n, got, n)
	}
}

func TestCounter_LargerN(t *testing.T) {
	const n = 2000
	got := Counter(n)
	if got != n {
		t.Fatalf("Counter(%d) = %d, want %d", n, got, n)
	}
}
