# Design Amazon's Sales Rank by Category

Design the feature that shows each product's sales rank within its
category — "#3 in Books > Science Fiction" — computed from the raw
firehose of completed orders.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: record sales, compute per-category rankings over a
  recent window, serve a product's rank (and a category's top-N).
  Out of scope? (Checkout itself, recommendations, real-time ranks.)
- The key scoping question: how fresh do ranks need to be? The answer
  decides whether this is a batch problem or a streaming problem.
- Put numbers on it: orders/month, order-log bytes/month, categories,
  products per category, rank-read QPS.

## Suggested defaults (if you want a starting point)

- 300 million orders per month; ranks refreshed hourly are acceptable
- Products belong to a category tree (rank in leaf and parent
  categories)
- Sales window for ranking: trailing 30 days
- Reading a rank must be fast (it renders on every product page)

## What good looks like

By the end you should have: estimates separating the huge write volume
(order log) from the modest serving data (ranks), a pipeline design
(order log → aggregate per category/product → sort → serve store), a
concrete MapReduce-style formulation with the key/value shapes at each
stage, a serving data model that answers both "rank of product X" and
"top N in category Y" cheaply, and a scaling story for hot categories
and incremental refresh.
