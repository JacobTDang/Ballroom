# Grading Rubric — Transcode Video: Long-Running Operations

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Operation resource model

- The start call returns an operation resource immediately (no held
  connections); the operation has a stable name (e.g.
  `operations/{id}`) and is GET-able like any other resource
  (AIP-151/121).
- The operation carries `done` plus a result that is exactly one of
  response or error once finished.

## 2. Progress & metadata

- In-flight state (progress percent, stage, ETA) lives in operation
  metadata, typed and documented — clearly separated from the
  terminal response.
- Metadata is advisory and can update on every poll without changing
  the operation's identity.

## 3. Terminal-state semantics

- Success, failure, and cancelled are distinct, detectable terminal
  states (error carries a canonical code; cancelled surfaces as
  CANCELLED, not a generic failure).
- Once done, the result is immutable — repolling never changes it.

## 4. Cancellation & idempotency of control

- Cancel is a method on the operation with stated best-effort
  semantics (the job may still finish).
- Repeated cancels are safe; cancel on an already-terminal operation
  is a defined no-op or defined error, not undefined.

## 5. Notification strategy

- A concrete polling contract: suggested interval, exponential
  backoff, and what a poll costs.
- The push alternative (webhook/queue) is named with its trade-off
  (delivery guarantees, endpoint management) — chosen or deferred
  with a reason.

## 6. Operation lifecycle & GC

- A stated retention window for finished operations and the behavior
  after it (NOT_FOUND, documented).
- A client that lost the operation id has a path (list operations
  filtered by video/requestor) or the loss cost is explicitly
  accepted.

## 7. Communication & trade-offs

- The sync-vs-async threshold is argued with numbers (why seconds of
  work could be synchronous but minutes cannot).
- Trade-offs stated for the major choices; the candidate drove the
  4 steps themselves.
