# Grading Rubric — Charge Exactly Once: Idempotency Keys

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. The double-charge scenario

- The timeout-after-commit ambiguity is articulated: after a timeout
  the client cannot distinguish a lost request from a lost response,
  so a blind retry risks a second charge.
- The mechanism is motivated by this scenario, not bolted on.

## 2. Key mechanics & replay

- Client-generated key (UUID or equivalent) sent with every create,
  with a stated scope — unique per merchant/account and method
  (AIP-155's request_id, Stripe's Idempotency-Key).
- The server stores key → outcome; a replay returns the *stored*
  response with its original status — success or failure alike —
  never a re-execution.

## 3. In-flight duplicates

- The second-request-while-the-first-executes case is decided, not
  ignored: wait for the first to finish, or reject with a retryable
  in-progress conflict — either accepted with the trade-off stated.

## 4. Payload-fingerprint conflict

- Same key + different body → an explicit error, never a silent
  replay or a second execution; a request fingerprint is stored
  alongside the key to detect it. Keys name an attempt, not a wish.

## 5. TTL & storage

- A retention window with a rationale tied to the retry tail, what
  expiry means for a very late retry, and a rough storage estimate at
  the stated volume.

## 6. Method coverage

- Which methods carry keys: non-idempotent creates and mutating
  custom methods. GET/PUT/DELETE identified as already idempotent by
  their standard-method semantics; the AIP-133 client-specified-ID
  alternative acknowledged as the other road to the same guarantee.

## 7. Communication & trade-offs

- Server-side dedup vs client-supplied resource IDs compared; failure
  modes (key-store outage, replay flag for observability) addressed
  honestly; the candidate drove the 4 steps themselves.
