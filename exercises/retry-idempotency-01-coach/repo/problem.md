# Charge Exactly Once: Idempotency Keys

You are designing the charge-creation endpoint for a checkout
platform: `POST /v1/charges` charges a customer's card. Clients call
it over flaky mobile networks and retry on timeout — and a retried
create must never charge the customer twice.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- The Charge resource: fields, its state machine (pending →
  succeeded / failed), and its resource name.
- Which operations are in scope? Create is the heart; get/list for
  reconciliation. Are refunds and captures in or out?
- Put numbers on it: charges per day, what fraction of requests time
  out client-side, and how long the retry tail runs.

## Suggested defaults (if you want a starting point)

- 5 million charges/day; ~1% of requests time out client-side and
  get retried
- Retries usually arrive within seconds, with a straggler tail up to
  24 hours
- One merchant API key per request; single region for now

## What good looks like

By the end you should have: the resource model and methods; the
timeout-after-commit ambiguity stated crisply (why the client cannot
know whether the charge happened); the idempotency-key mechanics end
to end — who mints the key, its scope, what the server stores, what a
replay returns, what happens when the same key arrives with a
different body, and what happens when it arrives while the first
attempt is still executing; a TTL with a rationale and a storage
estimate; and a clear call on which methods need keys at all.
AIP-155 and Stripe's Idempotency-Key are the canon to design against.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
