# Grading Rubric — Sales Rank by Category

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to rank computation + serving, with checkout and
  recommendations explicitly out.
- Freshness pinned down early (hourly batch acceptable) and used to
  justify a batch architecture.

## 2. Back-of-envelope estimates

- Orders/sec and monthly log size computed; the contrast drawn between
  the big input (order log) and the small output (ranks).
- Rank-read QPS recognized as the latency-critical number.

## 3. High-level design

- A pipeline: order events → durable log/object store → batch
  aggregation → sorted ranks → serving store → API; batch vs online
  boundaries marked.

## 4. The aggregation

- A MapReduce-style formulation with concrete key/value shapes:
  map to ((category, product), qty), reduce to window totals, sort
  within category to assign ranks.
- The trailing-window handled explicitly (daily partials combined, or
  recompute over 30 days — either, argued).

## 5. Category tree

- Sales roll up to ancestor categories deliberately (emit one pair per
  ancestor, or post-aggregate), with the cost acknowledged.

## 6. Serving data model

- A store answering both point reads ("rank of product X in category
  Y") and top-N lists without scanning — e.g. keyed rows plus a
  precomputed top-N per category.
- Rank reads never touch the batch pipeline.

## 7. Scaling story

- Aggregation sharded (by category hash) with hot categories/products
  addressed; top-N lists cached; incremental vs full recompute
  trade-off stated.

## 8. Communication & trade-offs

- Trade-offs stated (freshness vs cost, full vs incremental refresh),
  driven by the step-1 numbers; the candidate drove the 4-step
  structure.
