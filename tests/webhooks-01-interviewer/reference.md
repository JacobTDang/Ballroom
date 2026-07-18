# Reference Design — Order Event Webhooks

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
`subscriptions/{id}` (AIP-121): endpoint_url, event_types[], secret,
state (active | disabled), created_at. Ordinary CRUD — create, get,
list, patch, delete. Events: order.created/paid/shipped/refunded.
Numbers: 5M orders/day × 3 events ≈ 175 events/s average, ~1k/s peak;
200k endpoints, ~2% down at any moment — failure is the steady state,
not the exception.

## 2. Methods & wire shapes
Delivery is `POST {endpoint_url}` with a CloudEvents-shaped envelope:
`{id, type, time, source, api_version, data}`. Delivered = any 2xx
within a timeout (5 s); everything else — timeouts, 3xx, 4xx, 5xx —
is a failed attempt. One worked example payload included.

## 3. Delivery you can promise (the hard part)
- **At-least-once**, dedup by `id` on the merchant side — documented
  as a merchant obligation. Exactly-once refused: the 2xx can be lost
  after the merchant processed, so redelivery is unavoidable;
  pretending otherwise just hides the duplicate path.
- **Ordering**: none promised across events. `data` carries the
  order's monotonically increasing `sequence`; guidance: treat events
  as invalidation hints and GET current state when it matters.
- **Signing**: `X-Signature: t=<unix>, v1=HMAC_SHA256(secret,
  t + "." + body)`. Reject if |now − t| > 5 min (replay window).
  Rotation: subscription holds current + previous secret; sign with
  current, verify accepts either, retire after overlap — zero missed
  deliveries.
- **Retries**: backoff ≈ 1m, 5m, 30m, 2h, 8h, 24h, then dead-letter.
  Dead letters listable per subscription with manual redelivery.
  >24h of sustained failure → auto-disable + notify; re-enable by
  merchant action or a passing health probe.

## 4. Evolution & operations
Envelope versioned by `api_version`; evolution additive-only (new
event types are opt-in via selectors; new fields ignored by old
consumers — the versioning question's rules reapplied). Thin payload
chosen: ids + sequence + summary fields, merchants fetch the full
order — bounds staleness and payload drift at the cost of read-back
traffic (~1 extra GET per event, priced in). Isolation: per-endpoint
FIFO queues drained by a shared worker pool with a per-endpoint
in-flight cap of 1-2 — a hung endpoint blocks only its own queue.
Come-back herd: drain backlog at a per-endpoint rate cap (e.g.
10/s) so six hours of backlog doesn't crush the merchant that just
recovered. Push vs poll: webhooks for timeliness, plus a LIST events
API as the poll fallback — the dead-letter path and the poll API are
the same machinery.
