# Grading Rubric — Instagram

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped to photo posting + home feed (+ likes as a counter problem);
  explicit out-of-scope calls (stories, reels, DMs, explore).
- Estimates: posts/day, photo sizes across renditions, feed reads/day;
  the read-heavy ratio stated and used.

## 2. High-level design

- Media path (upload → object storage + CDN, multiple renditions)
  separated from the metadata/feed path (post metadata, follow graph,
  feed service).

## 3. Feed generation

- A concrete fanout decision: fan-out-on-write into per-user feed
  stores for normal accounts, with a deliberate celebrity answer
  (fan-out-on-read merged at request time, or argued alternative).
- Feed entries hold references (post IDs), media stays on the CDN.

## 4. Feed storage & pagination

- Materialized feeds in a fast store (e.g. Redis lists) with bounded
  length; cursor-based pagination; cold users rebuilt from the follow
  graph on demand.

## 5. Counters & interactions

- Likes/comments counts aggregated asynchronously or sharded-counter
  style — not a synchronous UPDATE on a hot row; eventual display
  consistency accepted and stated.

## 6. Scaling story

- Feed store and post metadata sharded by user; hot content absorbed
  by CDN + cache; fanout worker pool sized from the posts/day ×
  average-followers arithmetic.

## 7. Communication & trade-offs

- Trade-offs stated (write amplification vs read latency, freshness vs
  cost), driven by the estimates; the candidate drove the 4-step
  structure.
