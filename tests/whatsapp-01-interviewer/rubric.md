# Grading Rubric — WhatsApp

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped to 1:1 + group messaging with delivery states; explicit
  out-of-scope calls (calls, stories, payments).
- Estimates: concurrent connections, messages/sec, connections one
  chat server can hold — used to size the fleet.

## 2. High-level design

- Persistent connections (WebSocket/TCP) to chat servers; a session
  registry mapping user → server; message flow A → server A → registry
  lookup → server B → B narrated end to end.

## 3. Delivery semantics

- Sent/delivered/read receipts modeled as explicit acks flowing back
  through the same path.
- Offline recipients: messages queued server-side per device and
  drained on reconnect; at-least-once delivery with client-side dedup
  by message ID.

## 4. Group messaging

- Group send = fanout to member sessions via the registry; large-group
  cost acknowledged with a mitigation (fanout workers, message queue),
  not hand-waved.

## 5. Storage & encryption posture

- Messages stored (or deliberately not stored) with retention stated;
  end-to-end encryption acknowledged: servers route ciphertext, key
  exchange out of scope — and what that makes impossible server-side.

## 6. Scaling story

- Chat servers scaled horizontally with the registry (and its own
  availability) as the critical shared state; presence/last-seen
  handled without hammering the registry; multi-device mentioned.

## 7. Communication & trade-offs

- Trade-offs stated (delivery guarantees vs duplication, storage vs
  privacy), driven by the estimates; the candidate drove the 4-step
  structure.
