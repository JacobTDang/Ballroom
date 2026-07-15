# My Design: Sales Rank by Category

## Step 1 — Use cases, constraints, estimates

<!-- Orders/sec from the monthly figure; log bytes/month. Contrast the
     write volume with the small serving dataset (ranks). Freshness
     target and window. -->

## Step 2 — High-level design

<!-- Order events land in a log/object store -> batch aggregation ->
     sorted ranks written to a serving store -> API reads. Mark what
     is batch and what is online. -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - The aggregation as MapReduce-style steps: map order ->
       ((category, product), quantity); reduce to totals over the
       window; sort within category to assign ranks. Spell out the
       key/value shapes.
     - Category tree: how does a sale roll up to parent categories?
     - Serving data model: table/KV layout answering "rank of X" and
       "top N of Y" without scans.
     - Refresh: full recompute each hour vs incremental -- pick one and
       defend it. -->

## Step 4 — Scale it

<!-- Shard the aggregation by category hash; handle hot categories and
     hot products; cache top-N lists; keep rank reads off the batch
     path entirely. -->
