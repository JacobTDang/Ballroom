# Grading Rubric — Two Admins, One Settings Page

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Lost-update demonstration

- The failing interleaving is written out concretely (A reads v7,
  B reads v7, B writes, A writes → B's change gone) before any
  mechanism is proposed.

## 2. ETag mechanics

- Reads return an etag (header and/or resource field, AIP-154);
  writes carry If-Match (or the request etag field); a stale etag →
  412 / FAILED_PRECONDITION with a machine-readable reason.
- The etag's source is chosen and defended (server version counter
  vs content hash), and the etag is opaque to clients.

## 3. Strong vs weak validators

- The strong/weak distinction is explained (byte-equality guarantee
  vs semantic equivalence) and the choice for this API is justified.

## 4. Read-side conditionals

- If-None-Match → 304 on unchanged settings, with the payload/cost
  win quantified against the read rate (every console page load).

## 5. Client recovery protocol

- After a 412 the client's exact next steps are specified: re-GET,
  then rebase/merge or surface the conflict to the human — the
  recovery is part of the API's documented contract, not left to
  each client team.

## 6. Policy & defaults

- A deliberate call on unconditional writes: forbidden (write
  without If-Match → precondition-required error) or allowed —
  argued from the blast radius of a lost update.
- Locking is considered and rejected (or scoped) with reasons.

## 7. Communication & trade-offs

- Correctness vs client friction argued; where optimistic
  concurrency is overkill is named; the candidate drove the 4 steps
  themselves.
