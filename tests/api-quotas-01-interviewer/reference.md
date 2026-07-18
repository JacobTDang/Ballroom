# Reference Design — Rate Limits & Quotas

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & the two mechanisms
Two different promises: **quota** = the entitlement you sell
(100k req/month free, metered, resets on the billing cycle, overage
is a billing event) and **rate limit** = the protection you enforce
(60 or 1,000 req/min, absorbs bursts, guards the backend). Separate
windows, separate counters, separate errors — a monthly quota that
also throttles seconds is neither good billing nor good protection.

## 2. Wire shapes
Every response carries `X-RateLimit-Limit`, `X-RateLimit-Remaining`,
`X-RateLimit-Reset` for the tightest applicable limit. A limited call:
`429` + `Retry-After: <seconds>` + the platform-standard error body
(`status: RESOURCE_EXHAUSTED`, `reason: RATE_LIMITED` vs
`QUOTA_EXCEEDED`, plus the dimension that tripped) — clients branch
on fields, never on message strings.

## 3. Algorithm & placement (the hard part)
- **Token bucket** per key per endpoint-class: capacity = burst
  allowance (e.g. 2× the per-minute rate), steady refill = the tier
  rate. Bursty-but-honest clients pass; sustained abuse doesn't.
- The **fixed-window flaw** named: a 60/min fixed window admits 120
  requests in the two seconds straddling a boundary — why calendar
  windows only suit the billing quota, never the protective limit.
- **Cost weighting**: search debits 10 tokens, reads 1 — the bucket
  meters load, not request count.
- **Placement**: enforced at the 3 gateways. Local buckets with
  async counter sync (small over-admission, no added latency) over
  a synchronous shared store (exact, but a cross-node hop per
  request) — over-admission bounded by node count is the cheaper
  error. The quota counter, by contrast, is centrally metered
  (billing must be exact; seconds of lag are fine).

## 4. Operations & client experience
Documented client contract: on 429 honor Retry-After, else
exponential backoff with jitter; SDKs do this by default.
`GET /v1/limits` returns the caller's tiers, current usage, and
resets — discovery is an endpoint, not a support ticket. Fairness:
per-key buckets ahead of shared backend capacity; tiers are
parameters on the same mechanism. Prefer graceful degradation for
soft overage (queue-or-slow where the product allows) but a hard 429
for the protective limit — predictable beats clever. Limit changes
roll out flagged per key with before/after telemetry, because a
mis-set limit IS an outage for the key it hits.
