# My Design: Order Event Webhooks

## Step 1 — Scope & resource model

<!-- Subscriptions as resources (AIP-121): fields (endpoint URL,
     event-type selectors, secret), who CRUDs them, and the event
     catalog. Put numbers on orders/day, events/order, endpoints,
     and the down-at-any-moment fraction. -->

## Step 2 — Methods & wire shapes

<!-- The subscription CRUD surface, and the delivery itself: what an
     event POST to the merchant looks like — envelope fields, one
     worked example payload, and what counts as "delivered"
     (which response codes). -->

## Step 3 — The hard part: delivery you can promise

<!-- The guarantee: at-least-once + dedup by event id — and why
     exactly-once is refused. Ordering honesty. Signing: HMAC over
     payload + timestamp, the replay window, and rotating a secret
     without dropping deliveries. The retry schedule, dead-letter
     visibility, and endpoint auto-disable/re-enable. -->

## Step 4 — Evolution & operations

<!-- The envelope's version story and additive evolution rules;
     thin payload vs full resource; the fan-out arithmetic at your
     numbers, per-endpoint isolation so a dead merchant never delays
     a healthy one, and the thundering herd when a big endpoint
     comes back. -->
