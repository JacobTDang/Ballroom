# Error Surface for a Payments API

A payments platform is unifying its API errors. Today every endpoint
invents its own: some return `{"error": "oops"}`, some return HTTP 200
with `"status": "FAILED"`, one returns a stack trace. Clients retry
things they must never retry and give up on things that would have
succeeded. Design the error surface — the codes, the shapes, the
retry contract — for an API where clients act on failures with money
at stake.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Who consumes errors? SDKs branching on codes, merchant dashboards
  showing humans text, alerting pipelines counting classes of failure.
- What failure families exist for payments: validation, auth,
  insufficient funds, duplicate charge, downstream bank timeout,
  rate limits, internal bugs.
- What must a client be able to decide from an error alone?

## Suggested defaults (if you want a starting point)

- REST/JSON surface; SDKs in three languages you don't control
  release timing for.
- A charge attempt that times out is the scariest case: the client
  cannot tell whether money moved.
- Compliance forbids leaking internals or other tenants' existence.

## What good looks like

By the end you should have: a small fixed canonical code set mapped
to HTTP status codes (no per-endpoint inventions), a machine-readable
details structure (reason/domain/metadata) so clients branch on
fields rather than parsing messages, every code classified
retryable-or-not with Retry-After where it applies, the
developer-message vs end-user-text split, a stated no-leak policy
(including the 404-vs-403 disclosure call), and the shape multi-item
requests use to report per-item failures.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
