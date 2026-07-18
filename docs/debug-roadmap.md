# Debug Roadmap

A ladder of fourteen find-and-fix problems for getting interview-ready
at debugging, practiced in Ballroom's Debug category — each problem in
Python, Go, and C++. Every starter is plausibly-written code with one
real bug in it; the hidden tests reproduce the symptom, and your job is
the smallest correct fix. Debugging rounds are common in interviews and
badly under-practiced: the skill is a repeatable method, not luck.

**Cadence**: one rung per session, in order — the phases move from bugs
that announce themselves to bugs that only show up on the second call.
Repeating a rung in a second language afterward is a cheap, high-value
rep: the same logical bug wears each language's own disguise.

## How to practice (start here)

1. **Launch**: `ballroom` → Enter past the boot checks → `1` Practice →
   **Debug** → pick the next unsolved problem → pick a language.
2. **Read the problem pane** (`M-1`): it describes the *symptom* (a
   crash, a wrong answer, a hang-forever build time, a second call that
   lies) — never the cause. The bug is already in the editor.
3. **Reproduce first**: run the visible behavior in the terminal
   (`M-3`) before touching code. A bug you haven't reproduced is a bug
   you're guessing at.
4. **Talk to the tutor** (`M-2`) like a rubber duck with opinions:
   describe what you *expected* the line to do and what it actually
   does. Hints-first mode nudges before it names.
5. **Submit with `M-q`**. The hidden tests pin the symptom and the
   edge cases around it — a fix that dodges the test without fixing
   the cause won't survive them.

**The method** — for every problem, in order:

1. **Reproduce** — see the failure with your own eyes, smallest input.
2. **Localize** — bisect with prints/asserts until one line is guilty.
3. **Name the class** — off-by-one? aliasing? stale state? Saying the
   bug's name out loud is most of the fix.
4. **Minimal fix** — change what's wrong, nothing else. Don't refactor.
5. **Re-run** — the original repro *and* the neighbors (empty input,
   boundaries, calling it twice).

## The ladder

Phase 1 — Read the crash (the stack trace hands you the line):

- [ ] **Off-by-one** (`off-by-one-01`) — a loop that walks one step
      past the end; the crash names the exact index.
- [ ] **Guarded, but not enough** (`edge-guard-01`) — the early-return
      guard checks the wrong minimum; single-element input still walks
      off the end.

Phase 2 — Silent wrong answers (nothing crashes; the output lies):

- [ ] **The result that never escaped** (`shadow-var-01`) — an inner
      declaration shadows the outer result; the function returns its
      stale initial value.
- [ ] **Binary search boundary** (`bsearch-boundary-01`) — a
      lower-bound search whose initial bounds can't represent
      "past the end"; the last slot is unreachable.
- [ ] **The filter that skips its neighbor** (`iter-mutate-01`) —
      removing while iterating still advances the index; adjacent
      matches survive.
- [ ] **Wrong floor** (`base-case-01`) — a recursion whose base case
      returns the wrong value, collapsing every answer to zero.
- [ ] **Negative buckets** (`trunc-floor-01`) — timestamp bucketing
      rounds toward zero instead of down; pre-epoch times land in
      buckets that don't exist.
- [ ] **The bill that never settles** (`float-money-01`) — float
      equality on money; three dimes never equal thirty cents.
- [ ] **Backwards ties** (`sort-comparator-01`) — a multi-key
      comparator with one direction inverted, reached only on ties.
- [ ] **Dedupe by identity** (`identity-equality-01`) — equality
      checks compare *which object* instead of *what value*; equal
      records both survive.
- [ ] **The snapshot that edits the original** (`alias-copy-01`) — a
      "copy" that still shares storage with the live data; undo
      corrupts the thing it was saving.

Phase 3 — It's not wrong, it's slow:

- [ ] **The log builder that takes forever** (`quadratic-concat-01`) —
      string prepend in a loop; correct output, quadratic time. The
      hidden test has a stopwatch.

Phase 4 — Only fails the second time (stateful bugs):

- [ ] **The cache that confuses two questions** (`memo-key-01`) — a
      memo keyed on half the inputs; first call right, second call
      replays the wrong answer.
- [ ] **Right once, wrong forever** (`global-state-01`) — results
      accumulate in storage that outlives the call; every run after
      the first contains ghosts.

## Sibling roadmaps

The concurrency ladder's **deadlock-fix** (`deadlock-fix-01`) is this
ladder's rung fifteen in spirit — the same find-and-fix rep with
threads involved (`docs/concurrency-roadmap.md`). The implementation
ladder (`docs/implementation-roadmap.md`) is the inverse discipline:
building it right the first time.
