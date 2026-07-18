# Reference Design — Charge Exactly Once

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
`charges/{charge_id}` (AIP-122): amount, currency, source token,
state (pending → succeeded | failed), decline reason, created time.
5M charges/day ≈ 60/s steady; ~1% client-side timeouts means ~50k
replays/day the design must absorb without a second execution.

## 2. Methods & wire shapes
- `POST /v1/charges`, body = the resource (AIP-133); 201 with the
  created charge and its server-assigned name.
- A declined card is a *successful* API call returning
  `state: "failed"` + a machine-readable decline reason — not a 4xx:
  the request worked, the business outcome didn't.
- `GET /v1/charges/{id}` and paginated LIST for reconciliation.

## 3. Idempotency keys (the hard part)
- The ambiguity: a timeout after commit — the client saw nothing,
  the money moved. An unprotected retry double-charges.
- `Idempotency-Key: <uuid>` header, minted by the client once per
  logical attempt (AIP-155 / Stripe canon). Scope: unique per
  merchant account; the same key under another account is a
  different request.
- Server flow: atomically claim the key → execute → store
  (status, response body, request fingerprint) → respond. Replay:
  return the stored status + body verbatim, plus a replayed marker
  header for observability.
- In-flight duplicate: reject with a retryable 409
  `reason: IN_PROGRESS` and a small Retry-After — holding the second
  connection open through a payment is the worse trade.
- Fingerprint conflict: hash of the canonical body stored with the
  key; same key + different hash → 422 `KEY_REUSED` error. Never
  guess which request the client meant.

## 4. Evolution & operations
- TTL: 24h of guaranteed replay (covers the stated retry tail),
  documented loudly; a later retry is treated as a new request.
  ~5M keys/day × ~1 KB ≈ 5 GB/day in a KV store — trivial next to
  one double charge.
- Coverage: create and mutating custom methods (`:capture`,
  `:refund`) require keys; GET is safe, PUT/DELETE are idempotent by
  standard-method semantics (AIP-134/135) though keys are accepted
  uniformly for simpler client code.
- Alternative: client-specified charge IDs (AIP-133) give the same
  guarantee via ALREADY_EXISTS — clean, but pushes ID discipline onto
  every client; the header keeps resource naming server-owned.
