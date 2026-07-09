package main

import "time"

// RateLimiter allows at most `limit` calls per `window`.
type RateLimiter struct {
	limit  int
	window time.Duration

	windowStart time.Time
	count       int
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{limit: limit, window: window}
}

// Allow returns true if a new request should be allowed right now.
func (r *RateLimiter) Allow() bool {
	now := time.Now()
	if r.windowStart.IsZero() || now.Sub(r.windowStart) >= r.window {
		r.windowStart = now
		r.count = 0
	}
	if r.count >= r.limit {
		return false
	}
	r.count++
	return true
}
