# Reference Design — Mint.com

## 1. Use cases, constraints, estimates
Connect accounts, pull transactions, categorize, budgets + alerts.
10M users × 3 accounts × 30 txns/month ≈ 900M txns/month → **~350
writes/s** — modest, and daily-batch tolerant. Users read a few times a
month: **write-heavy**, the unusual shape worth naming. ~250 B/txn →
~2.7 TB/year.

## 2. High-level design
Ingestion pipeline: scheduler → **account-sync queue** → sync workers
(pull from institutions) → categorizer → transaction store → budget
checker → notification queue.
Read path: app → API → precomputed monthly rollups.

## 3. Core components
- **Sync:** one queued job per account daily; per-institution isolation
  (timeouts, retries, circuit breaker) so one slow bank can't stall the
  fleet. Credentials vaulted, never in app DBs.
- **Categorizer:** seller → category map (exact + heuristic), default
  bucket for unknowns, user overrides persisted and applied first.
- **Budgets/alerts:** category totals updated incrementally as txns
  land; alert fires once per crossing, delivered async, idempotent.
- **Data model:** append-heavy transactions table partitioned by month;
  monthly per-category rollup table serving all reads.

## 4. Scale
- Rollups make reads O(1); recompute via nightly batch (MapReduce
  shape: map txn → ((user, category, month), amount); reduce sums).
- Shard transactions by user; cold months to cheap storage.
- Notification fanout through its own queue + workers.
