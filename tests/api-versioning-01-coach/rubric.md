# Grading Rubric — Evolve the API Without Breaking Anyone

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Breaking-change taxonomy

- The concrete changes are sorted correctly (AIP-180): adding an
  optional field = safe; removing or renaming a field = breaking;
  tightening validation = breaking; changing a default or the meaning
  of a field = breaking.
- The enum-value addition gets its own treatment: safe only if
  clients are required to handle unknown values — a stated client
  obligation, not an assumption.

## 2. Wire-compatibility mechanics

- Clients must tolerate unknown fields; servers add but never remove
  or repurpose; a retired name/number is never reused with a new
  meaning.

## 3. Version surface

- Major version in the path (`/v1/`) per AIP-185; one major version
  per surface; minor/patch changes invisible on the wire.

## 4. Deprecation lifecycle

- Announce → dual-run window → sunset, with deprecation signaled in
  headers/docs and telemetry (per-client version usage) proving who
  is still on v1 before anything is turned off.

## 5. Stability tiers

- Alpha/beta/GA defined as contracts (AIP-181): what may break at
  each tier and with how much notice — so "it was beta" is a defined
  promise, not an excuse.

## 6. Migration story

- What a v1→v2 client migration actually touches; server-side v1
  shim over the v2 backend vs dual stacks, with the trade-off; the
  cost of the compatibility window stated.

## 7. Communication & trade-offs

- Velocity vs stability argued from the client mix (long-tail
  never-upgraders); the "run v1 forever" cost quantified rather than
  hand-waved; the candidate drove the 4 steps themselves.
