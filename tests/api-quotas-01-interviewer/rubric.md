# Grading Rubric — Rate Limits & Quotas for the Platform

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Quota vs rate limit distinction

- Long-window entitlement (per month/day, billable, tier-defining)
  and short-window protection (per second/minute, burst-absorbing)
  designed as two separate mechanisms with separate windows — not one
  blurred "limit".

## 2. Limit dimensions

- Dimensions (API key / user / IP / endpoint class) each justified by
  a concrete abuse case, not listed for completeness.
- Cost-weighted requests for expensive endpoints (the 10× search),
  so "requests per minute" measures load, not luck.

## 3. Algorithm choice

- Token bucket or sliding window chosen with the burst behavior
  explained (capacity vs refill).
- The fixed-window boundary flaw named explicitly (2× the limit
  straddling a window edge) — knowing why the naive thing fails.

## 4. The 429 contract

- 429 with a RESOURCE_EXHAUSTED-style machine-readable error body
  consistent with the platform error model, Retry-After honored,
  and X-RateLimit-Limit / -Remaining / -Reset headers specified
  (on successful responses too, not only on the 429).

## 5. Client guidance

- Documented client behavior: exponential backoff with jitter,
  SDK-level handling, and programmatic discovery of one's own
  limits/usage (an endpoint, not a support ticket).

## 6. Fairness & isolation

- One tenant's burst cannot starve others (per-key isolation ahead
  of shared capacity); paid tiers as different parameters, not
  different code paths; graceful degradation vs hard cutoff decided
  and argued.

## 7. Communication & trade-offs

- Protection vs developer experience argued; enforcement placement
  (edge vs service) with the distributed-counting trade-off
  (accuracy vs latency of shared counters) stated; the candidate
  drove the 4 steps themselves.
