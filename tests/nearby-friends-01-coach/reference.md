# Reference Design — Nearby Friends

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Opted-in friend pairs see each other within a radius, live as people
move; stranger discovery and navigation out. 100M location-sharing
users, ~10% concurrently active → 10M concurrent; ~30 s ping interval →
**10M / 30 ≈ 330K location writes/sec** — the defining number, and it's
write-heavy, unlike most of this track. One location record: lat/lng +
timestamp + user id, well under 100 B.

## 2. High-level design
Mobile client → location service (ingest the ping, write-heavy path) →
in-memory location store + geo index. A separate query path: for user
U, look up U's opted-in friends, check their current cells, return
matches within the radius.

## 3. Core components
- **Geo-index:** geohash buckets (or a quadtree) sized so one cell ≈
  the notification radius. Coarser cells mean fewer lookups but more
  false-candidate filtering; finer cells mean the opposite — pick one
  and say why. Two people meters apart can land in adjacent cells, so a
  lookup always scans the user's cell plus its neighbors, then filters
  by exact distance.
- **Push vs pull:** push wins here — pub/sub channel per user (or per
  cell) so a friend's move is delivered without the recipient polling.
  Fanout cost is per friend **pair**, not per user: a user with 200 friends
  sharing back potentially touches 200 channels per ping, so batch/
  coalesce updates rather than firing one message per ping per friend.
- **Privacy:** opt-in checked at read/notify time against the friend
  graph, not just hidden in the client UI — a stale or reshared client
  build must not leak location. Precision can be deliberately coarsened
  (round to a wider cell) before it ever leaves the query path.

## 4. Scale
- Freshness: every location record carries a short TTL; a user who
  stops pinging (offline, backgrounded) ages out instead of showing as
  a stale ghost. Ping frequency adapts — less often when the user's
  last few pings show they're stationary — cutting both battery and
  write load.
- Location store sharded by cell (co-locates the geo queries) or by
  user (co-locates a user's own writes) — state which and accept the
  trade-off. Dense urban cells get hot; handle it by finer-grained
  sub-cells in dense areas rather than a uniform grid everywhere.
- The friend-graph lookup (who are U's opted-in friends) is cached,
  since it changes far less often than location does; notification
  fanout scales as its own tier, separate from raw ingestion.
