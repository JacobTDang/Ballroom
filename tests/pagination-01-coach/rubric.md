# Grading Rubric — Paginate the Product Catalog

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Token design

- `next_page_token` is opaque and server-minted; clients never parse
  or construct it (AIP-158).
- Its contents are specified anyway: keyset position (last-seen sort
  key + id tiebreaker) plus a fingerprint of the query parameters.

## 2. Offset rejection with evidence

- OFFSET's failure is demonstrated, not asserted: O(n) skip cost at
  depth (page 1M scans and discards ~1M rows) and row drift under
  concurrent writes (skips/repeats when rows land before the cursor).

## 3. Stability under mutation

- Pinned semantics for inserts/deletes mid-walk: keyset guarantees no
  repeats of seen keyspace and no misses of rows that existed before
  the walk began; rows inserted behind the cursor are legitimately
  not seen.

## 4. page_size semantics

- Default and maximum stated; server may return fewer than requested
  without meaning "end"; end is signaled only by an empty
  `next_page_token`; zero/absent page_size falls to the default
  (AIP-158).

## 5. Token lifecycle & safety

- Expiry policy stated; garbage/tampered tokens → INVALID_ARGUMENT,
  never an empty success.
- A token is bound to its filter/order parameters — reuse after
  changing either is rejected via the fingerprint.

## 6. total_size / counting trade-off

- A deliberate call on offering counts: exact COUNT at 100M is
  priced (slow/expensive), alternatives named (estimate, omit,
  separate stats endpoint).

## 7. Communication & trade-offs

- Keyset vs offset vs cursor-table compared against the stated scale
  and client mix; the chosen orderings tied to real indexes.
