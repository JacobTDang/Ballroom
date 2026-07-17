# Sliding-Window Rate Limiter

The fixed-window limiter (rate-limiter-01) has a famous flaw: a burst
right before the boundary plus a burst right after it puts 2x the limit
through in one window's span. Fix it with a **sliding window**: allow a
request only if fewer than `limit` requests happened in the last
`window` milliseconds — measured from *this* request, not from a
boundary.

Time is injected (`nowMillis` is a parameter), so the tests are exact:
no sleeping, no flakiness — and no cheating with wall clocks.

## The invariant the tests enforce

- At most `limit` allowed requests in ANY `window`-millisecond span —
  including spans straddling old fixed-window boundaries.
- Requests age out precisely: an allowed request exactly `window` ms
  old no longer counts.
- Denied requests do not count against the window.

API: `NewSlidingWindow(limit int, windowMillis int64) *SlidingWindow`, `AllowAt(nowMillis int64) bool`. Timestamps are non-decreasing.

Think: what do you have to remember, and when can you forget it?
