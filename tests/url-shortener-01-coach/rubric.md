# Grading Rubric — Design Pastebin / Bit.ly

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Stated which use cases are in scope (shorten, redirect) and made a
  deliberate call on expiration, analytics, and custom aliases.
- Asked clarifying questions / stated assumptions rather than diving in.

## 2. Back-of-envelope estimates

- Write and read QPS derived from an assumed volume, with the
  read-heavy ratio stated (e.g. 10:1) and reflected in later choices.
- Storage sized over the link lifetime (bytes per record × volume),
  arithmetic shown.

## 3. High-level design

- End-to-end path for both flows: shorten (client → LB → write API →
  store) and redirect (client → LB → read API → store → 301/302).
- Components named and connected coherently; nothing load-bearing
  appears later that isn't on the diagram.

## 4. Short-code generation

- A concrete scheme: base-62 of an auto-increment ID, or hash
  (e.g. MD5 of URL + salt) truncated with collision handling —
  either accepted, but collisions must be addressed.
- Code length justified against the scale estimate (e.g. 62^7 ≈ 3.5
  trillion covers 3 years of writes with room).

## 5. Data model & storage choice

- Schema for the link mapping (short code, original URL, created/expiry,
  optional user/click count).
- SQL vs NoSQL choice justified by the access pattern (simple key
  lookup, huge read volume), not by fashion.

## 6. Redirect semantics

- Chose 301 vs 302 and explained the trade-off (caching at the browser
  vs analytics visibility).

## 7. Scaling story

- Caching for the read path (popular links, LRU/TTL, cache-aside or
  similar) — the single most important lever for a 10:1 read ratio.
- Database scaling beyond one box: replication for reads, sharding
  strategy for writes, and how the short-code scheme interacts with it.
- Hot-link handling (a viral link shouldn't take down a shard).

## 8. Communication & trade-offs

- Trade-offs stated for the major choices, driven by the step-1
  numbers; the candidate drove the structure (4 steps) themselves.
