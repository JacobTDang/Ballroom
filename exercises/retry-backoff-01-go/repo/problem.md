# Retry With Exponential Backoff

Call a flaky operation until it succeeds — waiting `base * 2^attempt`
between tries, capped at `cap`, giving up after `maxAttempts` with the
last error. The sleep function is injected, so the tests see exactly
what you'd have waited without waiting it.

The starter retries in a tight loop with a fixed delay — the thundering
herd in miniature.

## The invariant the tests enforce

- Success on attempt k stops immediately: exactly k calls, k-1 sleeps.
- The recorded delays are exactly `min(cap, base * 2^i)` for
  i = 0, 1, 2, ... — doubling, then flat at the cap.
- After `maxAttempts` failures: the operation's **last** error comes
  back, with exactly `maxAttempts` calls and `maxAttempts - 1` sleeps
  (no pointless sleep after the final failure).
- A first-try success sleeps zero times.

API: `Retry(op func() error, maxAttempts int, baseMillis, capMillis int64, sleep func(int64)) error`.
