# Reference Design — Notification System

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
One internal send API for push/SMS/email, called by other services; an
in-app inbox and campaign tooling out. 500M/day → **~5,000 sends/sec
average** (86,400 s/day rounds to 100K), peaking 5× at ~25K/sec.
Channels differ sharply: push is cheap and near-instant, SMS costs real
money per message and can be rate-limited by carriers, email is cheap
but slowest to land — each gets its own path, not one generic "send"
pipe internally.

## 2. High-level design
Caller service → ingest API (validates, assigns an idempotency key) →
preference/opt-out check → **per-channel queue** → channel-specific
worker pool → third-party provider (APNs/FCM, SMS gateway, email
service). Async end to end past the ingest API: the caller gets an
ack that the notification was accepted, not that it was delivered.

## 3. Core components
- **Preferences & capping:** opt-outs and per-channel preferences
  checked before enqueue, not after — a muted user's notification never
  reaches a queue. Per-user rate caps and quiet-hours windows enforced
  at the same check. Device tokens / phone numbers / addresses live in
  a registry, pruned when a provider reports a hard bounce or an
  invalid token.
- **Reliability:** every send carries an idempotency key derived from
  the triggering event, so a retried enqueue can't double-send.
  Provider failures retry with backoff; a message that exhausts
  retries lands in a dead-letter queue for inspection, not silently
  dropped. At-least-once delivery, chosen deliberately over
  exactly-once — the dedup key absorbs the duplicate risk cheaply,
  where exactly-once would need distributed transactions for no real
  gain here.
- **Prioritization:** priority classes (e.g. OTP/security = immediate
  lane, digests = bulk lane) as separate queues per channel, not one
  FIFO — a login code must not queue behind a million weekly digests.

## 4. Scale
- Each channel's worker pool scales independently — an SMS traffic
  spike doesn't starve push workers, because they're pulling from
  different queues entirely (bulkheading).
- Provider failover: a secondary SMS/email provider behind the same
  worker interface, swapped on sustained failure.
- Template rendering and delivery-tracking (sent/failed events for
  analytics) happen off the hot send path, feeding an events stream
  instead of blocking the worker that's trying to deliver the message.
