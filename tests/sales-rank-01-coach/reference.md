# Reference Design — Amazon Sales Rank by Category

## 1. Use cases, constraints, estimates
Record sales; hourly-refreshed per-category ranks over a trailing 30
days; fast rank reads (every product page). 300M orders/month →
**~120 orders/s** into a log; the servable output (ranks) is tiny by
comparison — the input/output size contrast is the shaping insight.

## 2. High-level design
Order events → durable log/object store → **hourly batch aggregation**
→ sorted ranks → serving store (KV/SQL) → rank API. Batch vs online
boundary drawn explicitly; reads never touch the pipeline.

## 3. Core components
- **Aggregation, MapReduce shape:**
  map: order → ((category, product), qty);
  reduce: sum over the window;
  sort within category by total → assign ranks.
  Trailing 30 days = combine 30 daily partials (cheap) rather than
  rescanning raw history.
- **Category tree:** emit one pair per ancestor category (cost
  acknowledged: ×depth), so a sale ranks in leaf and parents.
- **Serving model:** row per (category, product) → rank for point
  reads; precomputed top-N list per category for browse pages.

## 4. Scale
- Shard aggregation by hash(category); hot categories split further by
  product range.
- Cache top-N lists; rank point-reads behind the page cache.
- Hourly full recompute of window totals is simple and predictable;
  move to incremental daily-partial merging when cost demands.
