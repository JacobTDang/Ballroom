package main

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsUpToLimitWithinWindow(t *testing.T) {
	rl := NewRateLimiter(3, 10*time.Second)
	if !rl.Allow() {
		t.Fatal("expected 1st call allowed")
	}
	if !rl.Allow() {
		t.Fatal("expected 2nd call allowed")
	}
	if !rl.Allow() {
		t.Fatal("expected 3rd call allowed")
	}
	if rl.Allow() {
		t.Fatal("expected 4th call denied")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	if !rl.Allow() {
		t.Fatal("expected 1st call allowed")
	}
	if rl.Allow() {
		t.Fatal("expected 2nd call denied")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow() {
		t.Fatal("expected call allowed after window reset")
	}
}

func TestRateLimiter_MultipleWindowsInSequence(t *testing.T) {
	rl := NewRateLimiter(2, 40*time.Millisecond)
	if !rl.Allow() || !rl.Allow() {
		t.Fatal("expected first window's 2 calls allowed")
	}
	if rl.Allow() {
		t.Fatal("expected 3rd call in first window denied")
	}
	time.Sleep(50 * time.Millisecond)
	if !rl.Allow() || !rl.Allow() {
		t.Fatal("expected second window's 2 calls allowed")
	}
	if rl.Allow() {
		t.Fatal("expected 3rd call in second window denied")
	}
}

func TestRateLimiter_LimitOfOne(t *testing.T) {
	rl := NewRateLimiter(1, 10*time.Second)
	if !rl.Allow() {
		t.Fatal("expected 1st call allowed")
	}
	if rl.Allow() {
		t.Fatal("expected 2nd call denied")
	}
}
