# Reference Design — WhatsApp

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
1:1 + group messaging with sent/delivered/read states; calls, stories,
payments out. ~100M devices concurrently connected at peak; ~100B
messages/day → **~1.2M messages/sec average**, bursty around real-world
events. One chat server realistically holds on the order of **1M
persistent connections** (event-driven I/O, not one thread per
connection) → **~100 chat servers** cover peak concurrency alone,
before message throughput is even considered.

## 2. High-level design
Client ↔ persistent WebSocket/TCP connection ↔ one chat server. A
**session registry** (fast KV store) maps `user_id → chat_server_id`.
Send flow: A → server A → registry lookup for B → if B is connected,
forward to server B → push to B; if not, persist for later.

## 3. Core components
- **Delivery semantics:** sent/delivered/read are acks flowing the
  same path in reverse, updated on the sender's client. Offline
  recipients: the message is queued server-side (per device), drained
  on reconnect. At-least-once delivery — a redelivered message is
  expected, so the client dedups by message ID before displaying it.
- **Group messaging:** one send fans out to every member's current
  server via registry lookups; for large groups this is real fanout
  work, mitigated with a worker pool or an internal queue rather than
  the sender's connection blocking on N deliveries.
- **Storage & encryption:** messages persist only until delivered (or a
  short bounded window), not forever — state the retention choice.
  End-to-end encryption means the server routes ciphertext between
  connections; it cannot read content, search message bodies, or
  recover a lost message's plaintext. Key exchange is out of scope.

## 4. Scale
- Chat servers scale horizontally (stateless besides their live
  connections); the session registry is the shared state that must
  itself be replicated and available, since every message touches it.
- Presence/last-seen: batched or lazily-propagated updates instead of a
  registry write on every connect/disconnect blip — a chatty presence
  system would dwarf the message traffic itself.
- Multi-device: the registry maps a user to multiple active
  server/device entries, and a send fans out to all of them, not just
  one.
