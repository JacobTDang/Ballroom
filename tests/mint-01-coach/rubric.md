# Grading Rubric — Mint.com

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to account sync + categorization + budgets/alerts, with
  explicit out-of-scope calls (bill pay, investments, credit scores).
- Recognized the unusual shape: write-heavy, batch-tolerant, low read
  frequency — and said so.

## 2. Back-of-envelope estimates

- Transactions/sec derived from users × accounts × transactions/month —
  and the conclusion drawn that volume is modest and batch-friendly.
- Storage sized per transaction over a multi-year horizon.

## 3. High-level design

- Ingestion pipeline (institution pull → categorizer → store → budget
  check) drawn separately from the user-facing read path.
- Async boundaries marked: queues between pull, categorize, and notify.

## 4. Account sync

- Daily sync as queued jobs with a worker pool; slow or failing
  institutions isolated (timeouts, retries, per-institution queues)
  so one bad bank doesn't stall the fleet.
- Credential handling addressed (tokenized/vaulted, never in app DBs).

## 5. Categorization

- Seller → category mapping with a rules/lookup service, a default for
  unknown sellers, and user overrides that persist and retrain the
  mapping for that user.

## 6. Budgets & alerts

- Budget state updated incrementally as transactions land (or a clearly
  argued alternative), not recomputed by scanning history on read.
- Alert delivery is async (queue → notification service), idempotent,
  and doesn't fire repeatedly for the same crossing.

## 7. Scaling story

- Monthly spending overviews precomputed (rollup tables or MapReduce-
  style batch) so reads are cheap; transaction store sharded by user.
- Growth handled: old transactions cold-stored or partitioned by month.

## 8. Communication & trade-offs

- Trade-offs stated (freshness vs batch simplicity, write-time vs
  read-time computation), driven by the step-1 numbers; the candidate
  drove the 4-step structure.
