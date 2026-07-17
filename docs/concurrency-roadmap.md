# Concurrency Roadmap

A ladder of ten hands-on problems for getting interview-ready on
concurrency, practiced in Ballroom's Concurrency category — each
problem in Go, Python, and C++, with hidden tests that run under the
race detector (`go test -race`) and ThreadSanitizer
(`clang++ -fsanitize=thread`), so "it usually works" never passes.

**Cadence**: one problem per session, in order — each rung teaches a
primitive or pattern the next rung leans on. Go is the canonical
interview language here; doing a rung in a second language afterward is
a cheap, high-value rep.

## How to practice (start here)

1. **Launch**: `ballroom` → Enter past the boot checks → `1` Practice →
   **Concurrency** → pick the next unsolved problem → pick a language.
2. **Read the problem pane** (`M-1`): every problem states an invariant
   the tests enforce ("the count is exact", "consumers see every item
   exactly once", "shutdown never loses an accepted job"). Your job is
   to make the invariant hold under real parallelism — not to slow the
   code down until the race hides.
3. **Talk to the tutor** (`M-2`) like a colleague at a whiteboard: say
   which primitive you're reaching for and why. The hints-first mode
   nudges before it names.
4. **Submit with `M-q`**. The hidden tests run with the race detector /
   TSan where the language has one — a data race fails the run even
   when the answer happens to be right.

**The method** — for every problem, before typing:

1. Name the shared state.
2. Name the invariant the tests will check.
3. Choose the primitive (mutex, atomic, channel, condition variable,
   semaphore, once) and say *why* the cheaper one below it isn't enough.
4. Decide the shutdown story — who stops first, who waits, and how
   no work is lost. Most real-world concurrency bugs live here.

## The ladder

Phase 1 — Primitives (mutual exclusion & atomicity):

- [ ] **Racy counter** (`racy-counter-01`) — fix a data race with a
      mutex or atomic; the race detector is the teacher.
- [ ] **Bounded producer–consumer queue** (`bounded-queue-01`) — a
      fixed-capacity queue where producers block when full and
      consumers block when empty: condition variables (or a buffered
      channel) and the wake-up-loss trap.

Phase 2 — Coordination (many workers, one goal):

- [ ] **Worker pool** (`worker-pool-01`) — N workers draining a job
      queue, results collected without loss or duplication; clean
      close-of-work signaling.
- [ ] **Pipeline fan-out/fan-in** (`pipeline-01`) — a stage that fans
      work out to parallel workers and merges results, preserving
      nothing-dropped/nothing-duplicated through the merge.
- [ ] **Semaphore-bounded fetcher** (`bounded-fetch-01`) — run many
      tasks with at most K in flight at once; prove the bound is never
      exceeded (the tests watch the high-water mark).

Phase 3 — Patterns (state machines under contention):

- [ ] **Token-bucket rate limiter** (`token-bucket-01`) — a
      thread-safe Allow() with refill; correctness under concurrent
      callers, not just single-threaded arithmetic.
- [ ] **Once-only lazy init** (`lazy-init-01`) — expensive
      initialization that must happen exactly once no matter how many
      threads ask first (double-checked locking done right, or the
      language's once primitive).
- [ ] **Barrier rendezvous** (`barrier-01`) — a reusable barrier: N
      participants wait until everyone arrives, then all proceed;
      reuse across rounds is where naive versions break.

Phase 4 — Lifecycle & debugging (where production bugs live):

- [ ] **Graceful shutdown** (`graceful-shutdown-01`) — a worker system
      that, on stop, finishes accepted work, refuses new work, and
      never deadlocks or leaks a worker.
- [ ] **Fix the deadlock** (`deadlock-fix-01`) — debug-style: given
      code that deadlocks under contention (lock-ordering inversion),
      make the tests pass without serializing everything.

## Per-language notes

- **Go** runs with `-race`: any data race fails the suite outright.
  Prefer channels where they read naturally; use mutexes without shame
  where they don't.
- **C++** compiles with `-fsanitize=thread` (clang, C++17): TSan flags
  races and lock-order inversions at runtime. `std::mutex`,
  `std::condition_variable`, `std::atomic` are the whole toolbox here.
- **Python** has no race detector, and the GIL hides torn writes — but
  not lost updates, wake-up losses, or deadlocks. The tests assert
  observable behavior (exact counts, bounds, orderings) under real
  thread interleaving.

## After the ladder

Re-run rungs in your weaker languages; then the Implementation
category's systems half (bloom filter, consistent hashing, retry with
backoff — see docs/implementation-roadmap.md) pairs naturally: build
the component single-threaded there, then ask "what breaks under
concurrency?" here.
