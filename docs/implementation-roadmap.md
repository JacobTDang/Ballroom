# Implementation Roadmap

Eleven "build it from scratch" problems for the implement-X interview
round, practiced in Ballroom's Implementation category — each in Go,
Python, and C++ with real hidden tests. Two halves: **systems
components** (the pieces behind every system-design box you draw) and
**parsers & tooling** (the classic implement-a-thing screens).

**Cadence**: one problem per session. These are single-threaded on
purpose — build the component here, then the Concurrency ladder
(docs/concurrency-roadmap.md) asks what breaks when threads arrive.

## How to practice (start here)

1. **Launch**: `ballroom` → `1` Practice → **Implementation** → pick
   the next unsolved problem → pick a language.
2. **Read the invariants** in the problem pane (`M-1`) — every problem
   states exactly what the hidden tests will hold you to (fixed
   false-positive budget, remap bounds, precedence table, ...).
3. **State your plan to the tutor** (`M-2`) before coding: the data
   structure, the update rule, the edge cases you expect. Implement-X
   interviews grade the plan as much as the code.
4. **Submit with `M-q`**, then answer the complexity quiz honestly —
   these components are exactly the ones interviewers ask "and how
   does it scale?" about.

## Systems components

- [ ] **Fixed-window rate limiter** (`rate-limiter-01`) — the warm-up:
      windowed counting with a clock.
- [ ] **Sliding-window rate limiter** (`sliding-window-limiter-01`) — the
      fixed window's burst-at-the-boundary flaw, fixed; injectable
      time, precise eviction.
- [ ] **Bloom filter** (`bloom-filter-01`) — bit array + k hashes; the
      tests hold you to zero false negatives and a bounded
      false-positive rate.
- [ ] **Consistent-hash ring** (`consistent-hash-01`) — nodes on a
      ring with virtual nodes; adding/removing a node must remap only
      its neighborhood, and the tests measure exactly that.
- [ ] **LRU cache with TTL** (`ttl-cache-01`) — capacity eviction and
      time expiry interacting correctly (an expired entry is gone even
      if it was just used; expiry never revives evicted keys).
- [ ] **Retry with exponential backoff** (`retry-backoff-01`) — capped
      exponential delays with deterministic jitter, injectable sleep;
      give up honestly after the budget.

## Parsers & tooling

- [ ] **Event emitter** (`event-emitter-01`) — on/off/once/emit with
      exact ordering and safe removal-during-emit semantics.
- [ ] **Arithmetic tokenizer** (`tokenizer-01`) — numbers (with
      decimals), identifiers, operators, parens; positions tracked;
      garbage rejected loudly.
- [ ] **Glob matcher** (`glob-match-01`) — `*`, `?`, and `[a-z]`
      classes, built by hand (no regex delegation): the classic
      backtracking two-pointer.
- [ ] **INI parser** (`ini-parser-01`) — sections, key=value,
      comments, whitespace, later-key-wins; malformed lines are
      errors, not silence.
- [ ] **JSON subset parser** (`json-parser-01`) — objects, arrays,
      strings, integers, booleans, null; a real recursive-descent
      parser with position-carrying errors.

## Ground rules the tests enforce everywhere

- **No standard-library escape hatches** where the component IS the
  exercise (no `regexp` in the glob matcher, no `json` module in the
  parser, no `configparser`).
- **Loud failures**: malformed input errors with a reason; nothing
  silently returns a zero value.
- **Determinism**: anything time-based takes an injectable clock or
  sleeper — the hidden tests never sleep to make a point.

## After the ladder

Take any systems-half component back to the Concurrency category's
mindset: make your bloom filter or TTL cache thread-safe, then ask the
tutor to poke holes. And in System Design sessions, these are no longer
boxes — you've built the internals of the things you're drawing.

The API-design ladder (`docs/api-design-roadmap.md`) is the spec-side
counterpart: several of its questions pair one-to-one with problems
here — design the mechanism there, build it here.
