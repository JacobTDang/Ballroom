# Idempotency Key Store

A client retries the same request (same idempotency `key`) because the
first attempt's response got lost — not because it wants the operation
to happen twice. The store's job is to make retries safe: run the
operation once, and hand every retry the *original* result instead of
a second execution.

**This is not `ttl-cache-01` again.** The cache problem is about
eviction and capacity — what to throw away when the store is full. This
problem has no capacity limit at all; the interesting state lives
entirely in the *lifecycle of one request*: nothing seen yet, running
right now, or already finished — and a client that reuses your key for
a genuinely different request is a bug you must catch, not silently
overwrite.

A single request's life:

- **`Begin(key, fingerprint, now_ms)`** is called before the work
  starts. It returns one of:
  - **execute** — this key is new (or its record expired): go do the
    work, then call `Complete`.
  - **in-flight** — another call with this key is still running:
    don't do the work again.
  - **replay** — this key already completed: here is the *original*
    stored response, byte-identical, never recomputed.
  - If the key is live (in-flight or completed, not expired) but
    `fingerprint` doesn't match what was stored, that's a conflict —
    the caller is reusing a key for a different request body, which
    is always a loud error, never a silent replay or overwrite.
- **`Complete(key, response, now_ms)`** records the result for an
  in-flight request. Calling it for a key with no live in-flight
  record — unknown, already completed, or expired — is a loud error;
  there's nothing to complete.
- Every record (in-flight or completed) carries one deadline, `ttl_ms`
  out from whenever it was last touched (`Begin` starting it,
  `Complete` renewing it). Once `now_ms` reaches the deadline, the
  record is gone — `Begin` treats the key as brand new. This bounds
  both how long a stuck in-flight request is tracked and how long a
  completed response is replayable.

## The invariant the tests enforce

The full lifecycle matrix (execute → in-flight → replay, and expiry
resetting it to execute again); a replayed response is exactly what
`Complete` stored; a fingerprint mismatch on a live key is always
rejected, whether the key is in-flight or already completed; the
expiry deadline is exact; `Complete` on anything without a live
in-flight record is a loud error.

API: `IdempotencyStore(ttl_ms)`, `.begin_at(key, fingerprint, now_ms) -> (state, response)` where `state` is `"execute"`, `"in-flight"`, or `"replay"` (`response` is only meaningful for `"replay"`), `.complete_at(key, response, now_ms)`. Both raise `ValueError` on a fingerprint conflict or an invalid `Complete`.
