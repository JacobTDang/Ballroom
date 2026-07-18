# Design a Distributed Rate Limiter

Design a rate limiter for a public API platform: enforce per-client
request limits across a fleet of API servers, so no single client can
overwhelm the platform or another client's fair share.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use case: a limit enforced per client (API key? user? IP?)
  across many stateless API servers, not one box in isolation. Do
  different clients get different tiers? Per-endpoint limits on top of
  the per-client one?
- What happens when a client goes over — reject outright, and with
  what response? What does the client need to see to back off
  correctly?
- Put numbers on it: the decision QPS the limiter itself must sustain
  (every request needs a check, not just the ones that get limited),
  how many distinct clients exist and how many are active at once, and
  the memory footprint per client's state.

## Suggested defaults (if you want a starting point)

- ~50,000 requests/sec platform-wide need a rate-limit check
- 10 million registered API clients; ~500,000 active in any given minute
- Two tiers: free = 60 requests/min, paid = 6,000 requests/min
- Over the limit: `429 Too Many Requests` with a `Retry-After` header

## What good looks like

By the end you should have: stated assumptions with the decision-QPS
and memory arithmetic, a high-level diagram (limiter placed at the
gateway/middleware layer, backed by a counter store shared across the
fleet), a concrete algorithm choice defended against the alternatives
(and the fixed-window boundary-burst flaw named, not just avoided), how
the distributed check-then-increment race is made atomic, a deliberate
fail-open-vs-fail-closed call for when the counter store is down, and a
scaling story for the counter store itself (sharding, hot keys,
multi-region).

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
