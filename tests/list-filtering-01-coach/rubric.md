# Grading Rubric — Filter and Order the Fleet

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Filter surface

- A bounded grammar (AIP-160 flavor): named fields, a fixed operator
  set (=, !=, <, >, has, AND at minimum), concrete example strings —
  not "we accept SQL/arbitrary queries".
- Unfilterable fields rejected loudly at request time, never silently
  ignored.

## 2. order_by contract

- Whitelisted sortable fields with per-key direction
  (`order_by=last_seen_time desc, name`), multi-key allowed or
  explicitly not.
- A documented, stable default order (with an id tiebreaker).

## 3. Index honesty

- Promised filter/sort combinations mapped to actual indexes; the
  F×S combinatorial explosion named and bounded (composite indexes
  for the hot combinations, rejection or degraded path for the rest).

## 4. Pagination interaction

- Page tokens bound to their filter+order (fingerprint); changing
  either mid-walk invalidates the token with an explicit error, not
  silent weirdness (AIP-158).

## 5. Search vs filter separation

- Free-text/relevance search is a distinct parameter or endpoint
  backed by a search system, not bolted into the exact-match filter;
  the consistency trade-off (index lag) stated.

## 6. Validation & errors

- Malformed filter → INVALID_ARGUMENT with the offending position or
  field named, consistent with one error envelope.

## 7. Communication & trade-offs

- Expressiveness vs operability argued; scope cuts (OR? NOT? nested
  groups?) stated as decisions with reasons.
