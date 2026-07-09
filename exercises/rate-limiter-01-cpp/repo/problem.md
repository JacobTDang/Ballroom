# Fixed-Window Rate Limiter

Implement a rate limiter that allows at most `limit` calls within any
fixed window of time. When a call comes in:

- If the current window has expired (elapsed time since it started is
  `>= window`), start a new window and reset the count.
- If the count within the current window is already at `limit`, deny the
  call (`Allow()` returns `false`).
- Otherwise, count the call and allow it.

The very first call always starts the first window.

## Example

```
rl = RateLimiter(limit=3, window=10s)
rl.Allow()  // true  (1st in window)
rl.Allow()  // true  (2nd in window)
rl.Allow()  // true  (3rd in window)
rl.Allow()  // false (limit reached)
// ... 10s pass ...
rl.Allow()  // true  (new window)
```

## Constraints

- `limit >= 1`
- `window` is a positive duration.
- Calls arrive from a single goroutine/thread only — this is a
  *fixed-window* limiter, not a distributed/concurrent one (no locking
  required).
