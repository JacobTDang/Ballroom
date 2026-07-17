# Concurrent Token Bucket

A token-bucket rate limiter shared by many threads: `Allow()` takes a
token if one is available, `Refill(n)` adds tokens (an external ticker
calls it), and the bucket never holds more than its capacity.

The starter checks-then-decrements without synchronization — under
contention it hands out more tokens than exist.

## The invariant the tests enforce

- With capacity C and N > C concurrent `Allow()` calls, **exactly C**
  succeed — never C+1, never C-1.
- `Refill(n)` makes exactly `min(n, headroom)` more calls succeed —
  the bucket clamps at capacity.

API: `TokenBucket(capacity)`, `.allow() -> bool`, `.refill(n)`.

Think: "is a token available?" and "take it" must be one atomic step —
what makes them one step here?
