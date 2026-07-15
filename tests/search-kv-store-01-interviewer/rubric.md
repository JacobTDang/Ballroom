# Grading Rubric — Search Key-Value Cache

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to the cache (lookup, populate, evict, expire) with the search
  pipeline explicitly out of scope.
- Zipf/popularity skew named as the reason caching works at all.

## 2. Back-of-envelope estimates

- QPS from the monthly search figure; memory = entry size × working
  set; node count derived from per-node memory. Arithmetic shown.
- A deliberate call on entry size (full serialized page vs doc IDs)
  with the trade-off stated.

## 3. High-level design

- Read-through flow drawn: query API → cache → on miss, search backend
  → populate cache → respond.

## 4. LRU mechanics

- Hash map + doubly-linked list, with why each structure is needed
  (O(1) lookup and O(1) recency move/eviction).
- Eviction on insert-when-full described precisely (tail removal +
  map cleanup).

## 5. Distribution

- Consistent hashing of the normalized query to a cache node; minimal
  remapping on node join/leave stated as the reason.
- Key normalization addressed (case, whitespace, term order) as a hit-
  rate lever.

## 6. Freshness

- TTL and/or invalidation on index update, with the staleness budget
  from step 1 driving the choice.
- Cache-miss stampede on expiry addressed (request coalescing, jittered
  TTLs, or serve-stale-while-refresh).

## 7. Scaling story

- Hot-key handling (replicate hot entries or L1 per-API-node cache).
- Node failure = lost cache, not lost data; recovery behavior and the
  thundering herd on a cold node addressed.

## 8. Communication & trade-offs

- Trade-offs stated (memory vs hit rate, TTL vs invalidation
  complexity), driven by the step-1 numbers; the candidate drove the
  4-step structure.
