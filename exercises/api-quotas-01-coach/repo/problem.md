# Rate Limits & Quotas for the Platform

Your API platform is opening to third-party developers with free and
paid tiers. Design the limiting surface: what gets limited along which
dimensions, the algorithm and its burst behavior, exactly what a
limited client sees on the wire, and how one tenant's burst never
becomes everyone's outage.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Quota vs rate limit: which long-window entitlements are you selling,
  and which short-window protections are you enforcing? They are not
  the same mechanism.
- The dimensions: per API key? per user? per IP? per endpoint class?
  What abuse case motivates each?
- Are all requests equal — what does the expensive endpoint cost?

## Suggested defaults (if you want a starting point)

- 10,000 API keys; free tier 60 req/min and 100k req/month; paid
  tier 1,000 req/min
- The search endpoint costs roughly 10× a simple read
- Limits enforced at the edge; counters shared across 3 gateway nodes

## What good looks like

By the end you should have: the quota/rate-limit split stated as two
mechanisms with two windows; limit dimensions chosen against concrete
abuse cases, with cost weighting for expensive endpoints; an algorithm
choice with its burst behavior explained — and the fixed-window
boundary flaw named; the full 429 contract (error shape consistent
with your platform's error model, Retry-After, and the
limit/remaining/reset headers); documented client behavior
(backoff with jitter, programmatic limit discovery); and the fairness
story — tiering, isolation, graceful degradation vs hard cutoff, and
where the counters live given multiple gateways.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
