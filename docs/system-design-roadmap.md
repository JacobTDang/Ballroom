# System Design Roadmap

A 3+ month curriculum for getting interview-ready on system design, built
around [system-design-primer](https://github.com/donnemartin/system-design-primer)
as the reference text and Ballroom's System Design track for practice.

**Cadence**: ~2 study blocks + 1 practice session per week, alongside the
remaining NeetCode categories. A study block is reading + notes on one topic
cluster; a practice session is a Ballroom design session (coach or
interviewer).

## How to practice (start here)

You never need anything beyond this repo and `ballroom`. One session,
end to end:

1. **Launch**: run `ballroom` → Enter past the boot checks → `1` Practice
   → **System Design** → pick a question (start with *Design Pastebin /
   Bit.ly*).
2. **Pick the session style** — this is the "language" choice: **coach**
   (guided, 90 min — your first pass on every question) or
   **interviewer** (a timed 45-minute mock — your second pass, 1–2 weeks
   later; the picker marks questions ready for it with `mock due`).
3. **In the session**, three panes: the problem statement (top left,
   `M-1`), the tutor (`M-2`), and a terminal (`M-3`). Write your design
   into `solution.md` in the editor, working the 4 steps below. Talk to
   the tutor as you go — the coach teaches, the interviewer probes and
   watches the clock. Estimation help: `less ~/back-of-envelope.md` in
   the terminal pane.
4. **Submit with `M-q`** when done. Three things happen: the hidden
   grading rubric appears in your workspace, the tutor grades your
   `solution.md` against it dimension by dimension (recording the
   pass/fail — expect a strict grader; a fail with specific evidence is
   the useful outcome), and on the 8 solved questions a `reference.md`
   appears too — read it against your design before moving on.
5. **Between sessions**, check `ballroom` → Stats: **Rubric weak spots**
   ranks the dimensions you keep losing points on — when one keeps
   showing up, that's your next study block from Phase 1's list.

That loop — session → graded → weak spots → targeted study → next
session — is the whole roadmap in miniature. The phases below just
order the material.

**The method** — every practice session follows the primer's 4-step approach.
Internalize it until it's automatic:

1. **Outline use cases, constraints, and assumptions** — ask clarifying
   questions, state who the users are and what the system must do, then put
   numbers on it (users, QPS, storage, read/write ratio).
2. **Create a high-level design** — boxes and arrows covering the main
   components end to end.
3. **Design core components** — go deep where it matters: API shapes, data
   models, the hash/id scheme, the feed fan-out, whatever the question hinges
   on.
4. **Scale the design** — find the bottlenecks and fix them with the standard
   toolbox (load balancer, horizontal scaling, caching, sharding, queues),
   naming the trade-off every time.

---

## Phase 0 — Foundations (weeks 1–2)

Goal: speak the vocabulary and estimate confidently.

- [ ] Read the primer's "How to approach a system design interview question"
- [ ] Watch the Harvard scalability lecture; read the scalability article
      (clones → databases → caches → asynchronism)
- [ ] Performance vs scalability; latency vs throughput
- [ ] CAP theorem; consistency patterns (weak/eventual/strong); availability
      patterns (fail-over, replication); availability in numbers (99.9% vs
      99.99%)
- [ ] Back-of-envelope drills: memorize powers of two + latency numbers every
      programmer should know; do 3–4 estimation exercises (e.g. storage for
      Twitter, QPS for a URL shortener, bandwidth for video streaming).
      Cheat sheet: `docs/back-of-envelope.md` here, or
      `less ~/back-of-envelope.md` in any session's terminal pane.

## Phase 1 — Core building blocks (weeks 3–6)

Goal: know every box you'd draw and when to draw it. One block per study
session, with notes:

- [ ] DNS, CDN (push vs pull)
- [ ] Load balancers (L4 vs L7, HA pairs) and reverse proxies
- [ ] Application layer: microservices, service discovery
- [ ] Databases I — RDBMS scaling: replication (leader-follower,
      leader-leader), federation, sharding, denormalization, SQL tuning
- [ ] Databases II — NoSQL: key-value, document, wide-column, graph; and the
      SQL-vs-NoSQL decision drill (practice justifying the choice out loud)
- [ ] Caching: where (client, CDN, web server, database, application) and how
      (cache-aside, write-through, write-behind, refresh-ahead); cache
      invalidation pain
- [ ] Asynchronism: message queues, task queues, back pressure
- [ ] Communication: TCP vs UDP, REST vs RPC (and when GraphQL comes up)
- [ ] Security basics at interview depth

## Phase 2 — The eight solved questions (weeks 6–12)

Goal: practice the method on questions with known-good solutions. Each
question is done **twice** in Ballroom's System Design category:

1. First pass in **coach** style — guided walkthrough, then compare your
   design against the primer's solution and note what you missed.
2. Second pass (1–2 weeks later) in **interviewer** style — a timed mock,
   graded against the hidden rubric, no peeking.

Order (roughly increasing difficulty):

- [ ] Pastebin / Bit.ly (URL shortener) — the canonical starter
- [ ] Twitter timeline & search — fan-out on write vs read, feeds
- [ ] Web crawler — distributed work, dedup, politeness
- [ ] Mint.com — aggregation, batch jobs, budget alerts
- [ ] Data structures for a social network — graph storage, shortest paths
- [ ] Key-value store for a search engine — caching, LRU, consistent hashing
- [ ] Amazon sales ranking by category — MapReduce-style analytics
- [ ] Scaling a system to millions of users on AWS — the progressive scaling
      story, told end to end

## Phase 3 — OO design (interleaved with Phase 2)

Goal: the "design a class" round. These are code-shaped, so they live as
normal Ballroom coding exercises (real hidden tests) under the OO Design
category:

- [ ] LRU cache (already in Ballroom as the AI-assisted `lru-cache-01`)
- [ ] Hash map
- [ ] Call center
- [ ] Deck of cards
- [ ] Parking lot
- [ ] Chat server

## Phase 4 — Repetition & depth (months 3+)

Goal: consistency under pressure and depth where you're weak.

- Re-run Phase 2 questions in interviewer style until grades are consistently
  strong; track which rubric dimensions keep costing points.
- The seven that were once interviewer-only -- Dropbox, YouTube/Netflix,
  WhatsApp, Instagram, a rate limiter, a notification system,
  nearby-friends/Yelp -- now ship a coach variant and a reference design
  too, so they follow the same coach-then-mock rhythm as the rest.
- Deep dives driven by rubric weak spots — e.g. hand-waving cache
  invalidation twice in a row means a caching deep-dive study block.
- Read real-world architectures from the primer's company-blog list and map
  each back to the building blocks.

## Sibling roadmaps

The boxes you draw here have build-them-yourself counterparts:
`docs/implementation-roadmap.md` (bloom filters, consistent hashing,
rate limiters, caches — the components inside the boxes) and
`docs/concurrency-roadmap.md` (what breaks inside those components
when threads arrive). Interleaving them turns "I'd put a cache here"
into "I've built that cache."

The API-design ladder (`docs/api-design-roadmap.md`) is the contract
side of the same skill: this roadmap draws the boxes, that one
specifies exactly how clients talk to them.
