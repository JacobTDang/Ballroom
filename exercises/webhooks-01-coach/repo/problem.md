# Order Events: the Webhook Surface

Your e-commerce platform needs to push order events (created, paid,
shipped, refunded) to merchant-owned HTTPS endpoints. Merchants
subscribe, you deliver — including on the day a big merchant's
endpoint is down for six hours and comes back all at once.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- What is a subscription? Who creates it, what does it select
  (event types, all orders vs a store), what secrets does it hold?
- What delivery guarantee will you promise — and which one will you
  refuse to promise?
- Put numbers on it: orders/day, events per order, endpoints, how
  many endpoints are down at any moment, endpoint latency.

## Suggested defaults (if you want a starting point)

- 5 million orders/day, ~3 events per order over its lifetime
- 200,000 merchant endpoints; ~2% unreachable at any given time
- Median endpoint responds in 150 ms; p99 ~800 ms; some hang

## What good looks like

By the end you should have: subscriptions modeled as plain CRUD
resources (AIP-121 — this track's opening lesson reapplied); an honest
delivery guarantee (at-least-once, dedup by event id, and a sentence
on why exactly-once is refused); an honest ordering story; HMAC
signing with a timestamp and a replay window, plus secret rotation
that doesn't drop deliveries; a retry schedule with a dead-letter
story and endpoint auto-disable; a versioned event envelope
(CloudEvents-shaped is the canon) with the thin-payload vs
full-resource trade-off argued; and the fan-out arithmetic with
slow-consumer isolation so one dead endpoint never delays the rest.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
