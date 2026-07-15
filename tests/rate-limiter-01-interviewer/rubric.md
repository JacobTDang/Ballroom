# Grading Rubric — Distributed Rate Limiter

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped: per-client limits (and per-endpoint tiers) enforced across
  many API servers; what happens on limit (429 + Retry-After / headers)
  stated up front.
- Estimates: clients × state per client → memory footprint; decision
  QPS the limiter itself must sustain.

## 2. High-level design

- Limiter placed deliberately (gateway/middleware) with a shared
  counter store (e.g. Redis) behind the fleet; rules/config service
  separate from the hot decision path.

## 3. Algorithm choice

- Token bucket vs fixed window vs sliding-window log/counter compared,
  one chosen and defended (burst behavior, memory, accuracy).
- The boundary-burst flaw of fixed windows demonstrated, not just
  named.

## 4. Distributed correctness

- The check-then-set race identified; atomicity via Redis Lua script /
  atomic INCR+EXPIRE (or equivalent CAS) — not "we lock somewhere".
- Clock/skew and TTL cleanup of idle keys addressed.

## 5. Latency & availability

- Per-request overhead budgeted (one round trip to the store, pipelined
  or local-first); a two-tier option (local approximate + global
  authoritative) considered for scale.
- Fail-open vs fail-closed when the counter store is down: a deliberate
  choice with the consequence stated.

## 6. Scaling story

- Counter store sharded by client key; hot keys (one abusive client)
  handled; multi-region behavior stated (per-region limits or global
  sync trade-off).

## 7. Communication & trade-offs

- Trade-offs stated (accuracy vs latency vs memory), driven by the
  estimates; the candidate drove the 4-step structure.
