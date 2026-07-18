# Evolve the API Without Breaking Anyone

The library API is live: v1, public, and thousands of client apps you
do not control depend on it. Product now wants changes — some
cosmetic, some structural. Your job is the compatibility contract:
what may change in place, what forces a v2, and how a v2 rolls out
without stranding v1 clients.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- The change list on the table: renaming a field, adding an optional
  field, tightening a validation rule, adding a new enum value,
  restructuring a nested resource. Which are asks, which are needs?
- Who are the clients? SDK users vs raw HTTP; upgrade lag; whether
  any client can be forced to move.
- What does "broken" mean concretely — a 4xx that wasn't there
  before, a dropped field, a changed default?

## Suggested defaults (if you want a starting point)

- ~5,000 active client apps; median upgrade lag 9 months
- A long tail of clients that will effectively never upgrade
- One public surface, currently `/v1/`, no stability tiers yet

## What good looks like

By the end you should have: a correct breaking/non-breaking sorting of
the concrete change list (AIP-180 is the bar — and the enum-addition
case deserves its own sentence); the wire-compatibility rules both
sides must obey; where the version lives (AIP-185); an alpha/beta/GA
ladder (AIP-181) so stability is a defined contract; a deprecation
lifecycle with telemetry that proves who is still on v1; and the
migration story with its cost stated honestly.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
