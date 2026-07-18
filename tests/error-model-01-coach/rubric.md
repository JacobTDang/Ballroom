# Grading Rubric — Error Surface for a Payments API

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Canonical code discipline

- A small fixed code set (INVALID_ARGUMENT, UNAUTHENTICATED,
  PERMISSION_DENIED, NOT_FOUND, ALREADY_EXISTS, FAILED_PRECONDITION,
  RESOURCE_EXHAUSTED, UNAVAILABLE, INTERNAL — this flavor, AIP-193)
  each mapped to one HTTP status; no per-endpoint invented codes, no
  HTTP-200-with-failure-body.

## 2. Machine-readable details

- A structured details block (ErrorInfo-style `reason` + `domain` +
  `metadata`) so SDKs branch on fields; message strings explicitly
  non-contractual.

## 3. Retryability contract

- Every code classified retryable / not / conditional; Retry-After on
  RESOURCE_EXHAUSTED and UNAVAILABLE; the timeout-after-charge
  ambiguity named, with the client guidance pointing at idempotent
  retry (the mechanism is a later question — the contract belongs
  here).

## 4. Human vs machine split

- Developer-facing `message` distinct from end-user-safe text (a
  localized display string or a documented mapping); neither doubles
  as the branching surface.

## 5. No-leak discipline

- Stack traces, table/host names, and other tenants' existence never
  appear; the 404-vs-403 disclosure policy is an explicit decision
  (e.g. 404 for resources the caller can't know exist).

## 6. Partial failure shape

- Multi-item requests report per-item status in request order,
  each slot a result-or-error of the same canonical shape — no
  silent drops, no all-or-nothing pretense.

## 7. Communication & trade-offs

- Precision vs disclosure vs compatibility argued; error surface
  treated as evolvable API (additive codes/details only), not
  freeform text.
