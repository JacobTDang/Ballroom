# Reference Design — Evolving Without Breaking

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & clients
~5k client apps, median 9-month upgrade lag, a never-upgrading tail.
"Breaking" defined observably: a request that succeeded yesterday
fails today, or a response field a client read yesterday is gone or
means something else.

## 2. The taxonomy (AIP-180)
- Add optional field → **safe** (old clients ignore it — because the
  contract obliges them to tolerate unknown fields).
- Remove or rename a field → **breaking** (a read that returned data
  now returns nothing; rename = remove + add).
- Tighten validation → **breaking** (yesterday's valid request now
  400s).
- Change a default or a field's semantics → **breaking**, the sneaky
  kind: the wire shape is identical and behavior still changes.
- Add an enum value → **subtle**: safe only because clients are
  *required* to treat unknown enum values as a defined fallback;
  without that stated obligation it is breaking.

## 3. The compatibility contract
Clients: tolerate unknown fields and unknown enum values. Servers:
additive-only within a major; never reuse a retired name with new
meaning. Version surface (AIP-185): major in the path — `/v1/` —
one live major per surface plus at most one in deprecation;
minor/patch invisible on the wire. Stability tiers (AIP-181):
alpha (may break anytime, allowlisted), beta (breaking changes with
notice), GA (no breaking changes, period).

## 4. Rollout & operations
v2 ships as `/v2/` beside `/v1/`; the v1 surface becomes a thin
translation shim over the v2 backend (one implementation, two
faces) — dual stacks drift, shims don't. Deprecation lifecycle:
announce with a date, emit `Deprecation` + `Sunset` headers, publish
per-client usage telemetry internally, chase the top offenders,
brown-out before turn-off. Budget honestly: the shim is code you
maintain for years — the price of 5k clients you don't control; the
alternative (breaking them) is priced in support tickets and churn,
which is why in-place breaking changes are never the cheap option.
