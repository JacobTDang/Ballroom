# Grading Rubric — Nearby Friends

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped: opted-in friends see each other within a radius, updating
  as people move; explicit out-of-scope calls (strangers/discovery,
  navigation).
- Estimates: active users × ping interval → location writes/sec (the
  defining number, and it's write-heavy); state per user.

## 2. High-level design

- Mobile clients ping location periodically → location service →
  in-memory location store + geo index; a query/notify path that tells
  a user which friends are near.

## 3. Geo-indexing

- A concrete spatial index: geohash buckets or quadtree, with the
  cell-size trade-off (precision vs candidates scanned) and the
  neighbor-cell problem at boundaries handled (search adjacent cells).

## 4. Push vs pull

- A deliberate delivery model: pub/sub channels per user/cell pushing
  updates to nearby friends' devices, or periodic pull — chosen and
  defended with the friend-pair fanout cost worked through.

## 5. Freshness & lifecycle

- Locations carry TTLs; stale/offline users age out rather than
  showing ghosts; update frequency adapted (moving vs stationary) to
  save battery and write volume.

## 6. Privacy

- Opt-in per friend pair enforced at read time (not just UI);
  precision reduced where appropriate; location history deliberately
  not retained (or retention stated).

## 7. Scaling story

- Location store sharded (by user or by cell) with hot urban cells
  handled; the friend-graph lookup cached; websocket/notification
  fanout scaled separately from ingestion.

## 8. Communication & trade-offs

- Trade-offs stated (freshness vs battery/write load, precision vs
  privacy), driven by the estimates; the candidate drove the 4-step
  structure.
