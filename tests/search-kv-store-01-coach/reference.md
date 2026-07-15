# Reference Design — Search Key-Value Cache

## 1. Use cases, constraints, estimates
10B searches/month → **~4K QPS**. Zipf skew is why caching works: a
small fraction of distinct queries is most of the traffic. Entry =
normalized query → serialized results (~250 KB full page, or ~1 KB of
doc IDs — a deliberate call). At 250 KB × 40M hot entries ≈ 10 TB →
**a cluster**, not a box; doc-ID entries fit far smaller. ~1 hr
staleness budget. Hit latency: single-digit ms.

## 2. High-level design
Query API → cache layer (read-through): hit → return; miss → search
backend → populate → return.

## 3. Core components
- **LRU mechanics:** hash map (O(1) lookup) + doubly-linked list (O(1)
  move-to-front and tail eviction) — why both structures is the core
  data-structure question.
- **Placement:** consistent hashing of the **normalized** query →
  node; minimal remapping on membership change.
- **Normalization:** case, whitespace, term order — directly buys hit
  rate.
- **Freshness:** TTL ≈ the staleness budget; on index update, versioned
  keys or explicit invalidation. Expiry stampedes: jittered TTLs or
  serve-stale-while-refresh.

## 4. Scale
- Hot keys (breaking news): small L1 cache in each API node, or
  replicate hot entries across N nodes.
- Node loss = lost cache only; cold-node thundering herd absorbed by
  request coalescing.
- Memory pressure: shrink entries (doc IDs, compression) before adding
  machines.
