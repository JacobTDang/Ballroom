# Reference Design — Fleet Filtering & Ordering

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope
Dashboards filter by status/site/firmware/last_seen/labels and sort by
name or recency; the 200k-device customer defines the worst case.
Search-by-name is a separate need with separate machinery.

## 2. Wire contract
`GET /v1/devices?filter=...&order_by=...&page_size=...&page_token=...`
Grammar (AIP-160 subset): `field op value` with ops `= != < > <= >=`,
`labels:key` for has-label, joined by `AND` only.
Examples:
- `filter=status = OFFLINE AND site = "sites/berlin-3"`
- `filter=firmware_version < "2.4.0" AND last_seen_time > "2026-07-01T00:00:00Z"`
- `order_by=last_seen_time desc, name` (default: `name, id` — stable
  via the id tiebreaker).
Unknown field or operator → `400 INVALID_ARGUMENT` naming it.

## 3. The hard part — promises vs indexes
Promise only what's indexed: composite indexes for the shipped
combinations (status+site, firmware+last_seen, each × the two sort
keys). The full F×S explosion (dozens of composites) is named and
refused: unindexed combinations return
`INVALID_ARGUMENT: unsupported filter/order combination` rather than
timing out at 200k rows. OR and NOT are cut in v1 (each multiplies
the index matrix); labels are served by an inverted index.
Page tokens carry a fingerprint of (filter, order_by); any change →
`INVALID_ARGUMENT: token does not match request` (AIP-158).
Free-text search is `q=` backed by a search index (relevance-ranked,
seconds-stale, no token compatibility with filter walks) — declared
eventually consistent, unlike filters which read the primary.

## 4. Evolution & operations
A field becomes filterable only after its index is built and
backfilled — the API doc is the contract, so it ships last, not
first (AIP-180: adding is safe, un-promising is breaking). Query cost
guard: per-request row-scan budget; over-budget → error advising a
tighter filter. Malformed-filter errors carry
`{ field, position, reason }` in details.
