# Reference Design — Social-Network Data Structures

## 1. Use cases, constraints, estimates
Friend lookups (fast) + shortest path between two users (~1s ok). 100M
users × 1K friends ≈ **100B edges**; adjacency lists at ~8 B/edge ≈
800 GB + overhead → **must be partitioned**. That constraint drives
everything.

## 2. High-level design
Client → API → **lookup service** (user → person-server) → person
servers, each holding a shard of adjacency lists in memory. Path
queries run in the API layer, calling person servers per hop.

## 3. Core components
- **Data model:** adjacency list per user, sorted array of user ids on
  its person server.
- **Sharding:** hash(user id) → server; the lookup service is a tiny
  replicated map (config service), cacheable client-side.
- **Distributed BFS:** frontier grouped by owning server each hop → one
  **batched** RPC per server per hop; shared visited set; depth cap
  (~6). Naive per-node RPCs is the failure mode to name.
- **Bidirectional BFS:** expand from both ends, meet in the middle —
  turns b^d into 2·b^(d/2); the expected latency fix.
- **Writes:** friend/unfriend touches two adjacency lists on (usually)
  two servers — transactional pair-write or async with reconciliation;
  say which and the consistency consequence.

## 4. Scale
- Hot users' lists cached (LRU); person servers replicated for reads.
- Heavy tail: multi-million-degree accounts chunked/paginated.
- Lookup service can't bottleneck: replicate + client caching.
