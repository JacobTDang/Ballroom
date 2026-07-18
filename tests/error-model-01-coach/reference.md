# Reference Design — Payments Error Surface

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope
Three consumers: SDKs (branch on codes/reasons), dashboards (show
humans safe text), alerting (count classes). A client must decide from
the error alone: is this my bug, may I retry, when, and did money
move.

## 2. The canonical set and envelope (AIP-193)
One envelope on every non-2xx response:
```json
{
  "error": {
    "code": "FAILED_PRECONDITION",
    "message": "card_expired: card ending 4242 expired 2026-05",
    "details": [{
      "reason": "CARD_EXPIRED",
      "domain": "payments.example.com",
      "metadata": { "card_last4": "4242" }
    }]
  }
}
```
Codes → HTTP: INVALID_ARGUMENT 400, UNAUTHENTICATED 401,
PERMISSION_DENIED 403, NOT_FOUND 404, ALREADY_EXISTS 409,
FAILED_PRECONDITION 422, RESOURCE_EXHAUSTED 429, INTERNAL 500,
UNAVAILABLE 503. Insufficient funds = FAILED_PRECONDITION +
reason INSUFFICIENT_FUNDS (a domain condition, not a malformed
request). Duplicate charge = ALREADY_EXISTS. Never 200-with-failure.

## 3. The retry contract
- Retryable as-is: UNAVAILABLE, INTERNAL (bounded, with backoff);
  RESOURCE_EXHAUSTED after Retry-After.
- Never retry unchanged: INVALID_ARGUMENT, PERMISSION_DENIED,
  FAILED_PRECONDITION (fix the condition first).
- The scary case — timeout on a charge: the client cannot know if
  money moved, so the documented contract is "retry with the same
  idempotency key" (mechanism specified in the idempotency question;
  the error model's job is saying so in the charge endpoint's docs
  and NEVER returning an ambiguous non-canonical shape).
- SDKs branch on `code` + `details[].reason` only; `message` is
  explicitly non-contractual and may change any release (AIP-180).

## 4. Evolution & operations
Additive-only: new reasons and metadata keys are safe; removing or
re-meaning either is breaking. No-leak: internals (stack, hosts,
queries) never serialize; cross-tenant probes get 404 (existence is
the secret — 403 would confirm it) — stated as policy, applied
uniformly. Batch endpoints return per-item
`{ "index": i, "payment": {...} | "error": {...} }` in request order.
End-user text lives in a `localized_message` detail with locale, kept
apart from the developer message.
