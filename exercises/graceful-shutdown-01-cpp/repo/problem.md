# Graceful Shutdown

A little job server: `Submit(v)` hands a job to a pool of workers,
`Stop()` shuts the system down **gracefully** — every job accepted
before Stop gets handled, new submissions are refused, and Stop only
returns when the workers have drained everything and exited.

The starter's Stop flips a flag and returns immediately: queued jobs
are abandoned mid-flight. This is where most real-world concurrency
bugs live — not in the hot path, in the shutdown.

## The invariant the tests enforce

- Every job accepted (`Submit` returned true) is handled, even when
  `Stop()` is called with a full queue.
- After `Stop()` returns: `Submit` returns false, and nothing new gets
  handled.
- `Stop()` itself completes (no deadlock, bounded time), and calling it
  is safe exactly once from anywhere.

API: `class Server { Server(int workers, std::function<void(int)> handle); bool Submit(int v); void Stop(); ~Server(); }`. Compiled with `-fsanitize=thread`.

Think: what tells the workers "no more work is coming", and who waits
for whom?
