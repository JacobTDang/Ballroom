# Reference Design — Distributed Rate Limiter

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Per-client limits (API key), plus per-endpoint tiers, enforced across a
stateless API fleet; over limit → **429 + `Retry-After`**, out of scope.
~50K req/s platform-wide need a check — that's the limiter's own load,
independent of any one client's limit. 10M clients, **~500K active/min**;
counter state ~16–32 B/client → **under 20 MB** resident, trivial for
one Redis primary.

## 2. High-level design
Client → LB → API gateway (limiter runs here, before the request
reaches app servers) → shared counter store (Redis) for the
increment-and-check. Rules/tiers live in a small config service, cached
at the gateway, not fetched per request.

## 3. Core components
- **Algorithm:** token bucket — smooth refill, allows a bounded burst,
  O(1) state (tokens, last-refill-ts) per client. Fixed window rejected
  by demonstration: a client can send its full quota right before a
  window resets and its full quota right after — 2× the limit in a
  moment, right at the boundary. Sliding-window log is the accurate
  alternative but costs O(requests) memory per client; a sliding-window
  counter approximates it in O(1) — reasonable substitute, state the
  accuracy trade-off.
- **Distributed correctness:** the read-modify-write (check tokens,
  decrement, maybe refill) is race-prone across concurrent requests
  from the same client hitting different gateway nodes. Fix: a single
  Redis Lua script does refill + check + decrement atomically in one
  round trip — not "check in app code, then write." `EXPIRE` on the
  key so idle clients' state ages out instead of growing forever.
- **Config propagation:** tiers pushed to gateways via a config service
  + local cache with a short TTL, not read from Redis per request.

## 4. Scale
- Latency: one Redis round trip per request, sub-ms same-AZ; pipeline
  where possible. A local, approximate token count (decremented
  in-process, resynced periodically) trades some accuracy for zero
  network hop on the hot path at very high QPS.
- Fail-open vs fail-closed when Redis is unreachable: fail-open (allow
  the request) protects availability but exposes the platform during
  an outage; fail-closed protects the platform but takes down every
  client's traffic with one dependency. Pick one and say which failure
  mode you're optimizing against.
- Counter store sharded by client key (consistent hashing); one abusive
  client is one hot key on one shard — bound the blast radius, don't
  reshard the world for it.
- Multi-region: per-region limits enforced independently (simple, can
  over-admit a client that spreads across regions) vs a global sync
  (accurate, adds cross-region latency to the hot path) — a deliberate
  trade-off, not an afterthought.
