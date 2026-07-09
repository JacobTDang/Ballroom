# Racy Counter

`Counter` (Go/C++) / `run` (Python) spins up `n` concurrent workers, each
incrementing a single shared counter exactly once, then returns the final
count. It's supposed to always return exactly `n`.

It currently has a **data race**: the increment isn't atomic, so two
workers can read the same value before either writes back, silently
losing an increment under contention.

## Task

Find and fix the race. The fix should not serialize all the work onto a
single goroutine/thread — synchronize the shared state instead (an
atomic operation or a lock around the read-modify-write).

Go tests run with `-race`; C++ tests are built with ThreadSanitizer
(`-fsanitize=thread`) — both will report exactly *where* the race is
happening if you want a hint before asking the tutor.

## Constraints

- `n >= 1`
