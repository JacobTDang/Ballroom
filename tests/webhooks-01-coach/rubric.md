# Grading Rubric — Order Events: the Webhook Surface

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Subscription model

- Subscriptions are plain CRUD resources (AIP-121): endpoint URL,
  event-type selectors, a signing secret; created/listed/patched/
  deleted like any resource — not a bespoke registration RPC.

## 2. Delivery semantics

- At-least-once declared honestly, with consumer-side dedup keyed on
  a unique event id.
- Exactly-once explicitly refused, with the reason (the ack can
  always be lost after processing).

## 3. Ordering honesty

- No global ordering promise. Either per-resource sequence numbers or
  explicit "treat events as hints — fetch current state on receipt"
  guidance; silent implied ordering is a miss.

## 4. Security

- HMAC signature over payload + timestamp with a bounded replay
  window; verification steps stated from the merchant's side.
- Secret rotation without downtime (overlapping dual secrets or
  versioned signatures).

## 5. Retry & dead-letter policy

- An exponential backoff schedule with a terminal give-up; failed
  deliveries visible to the merchant (dead-letter list or dashboard),
  not silently dropped.
- Persistent failures auto-disable the endpoint, with a re-enable
  path (manual or probe-based).

## 6. Event schema & evolution

- A versioned envelope (CloudEvents-shaped or equivalent: id, type,
  time, source, data) with additive-only evolution rules.
- Thin payload vs full resource argued (staleness vs fetch-back
  traffic), a choice made.

## 7. Scale & isolation

- The fan-out arithmetic done at the stated numbers (events/sec in
  and deliveries/sec out).
- Slow-consumer isolation: per-endpoint queues/workers so one dead
  merchant never delays healthy ones; the come-back thundering herd
  (hours of backlog at once) addressed.

## 8. Communication & trade-offs

- Push vs poll vs streaming compared against client capabilities;
  trade-offs tied to the step-1 numbers; the candidate drove the 4
  steps themselves.
