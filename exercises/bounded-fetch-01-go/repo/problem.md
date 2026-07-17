# Semaphore-Bounded Parallel Runner

Run a batch of tasks concurrently — but never more than `limit` at the
same time. Think "parallel fetcher that mustn't hammer the upstream".

The starter launches everything at once: fast, and exactly what takes
the downstream service off the air.

## The invariant the tests enforce

- Every task runs exactly once.
- The number of tasks in flight never exceeds `limit` (the tasks
  themselves measure the high-water mark).
- With `limit` > 1 and slow tasks, real parallelism happens (high-water
  of at least 2), and `limit` = 1 degrades to strictly serial.

API: `RunLimited(tasks []func(), limit int)`. Tests run with `-race`.

Think: what *is* a semaphore in this language, and where exactly must
the acquire/release go so the bound covers the task body — not just its
launch?
