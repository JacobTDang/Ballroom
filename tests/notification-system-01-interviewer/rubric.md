# Grading Rubric — Notification System

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped: one internal API through which services send push/SMS/email;
  explicit out-of-scope calls (in-app inbox, marketing campaign
  tooling).
- Estimates: notifications/day → peak sends/sec; per-channel cost and
  latency differences acknowledged.

## 2. High-level design

- Ingest API → validation/preference check → per-channel queues →
  channel workers → third-party providers (APNs/FCM, SMS gateway,
  email service); async end to end, with the reason stated.

## 3. User preferences & capping

- Opt-outs and per-channel preferences checked before enqueue; per-user
  rate caps / quiet hours to prevent spam; device-token / address
  registry maintained (and pruned on bounces).

## 4. Reliability

- Idempotency keys so one triggering event can't double-send through
  retries; retries with backoff on provider failure; dead-letter queue
  with what happens to it.
- At-least-once delivery chosen and defended (vs the cost of exactly-
  once).

## 5. Prioritization & isolation

- Priority classes (OTP now, digest later) as separate queues or
  priority lanes; one slow provider or channel can't starve the others
  (bulkheading per channel).

## 6. Scaling story

- Workers scale per channel independently; provider failover; template
  rendering and delivery tracking (sent/failed events for analytics)
  without blocking the send path.

## 7. Communication & trade-offs

- Trade-offs stated (latency vs batching, duplication vs loss), driven
  by the estimates; the candidate drove the 4-step structure.
