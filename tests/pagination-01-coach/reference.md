# Reference Design — Catalog Pagination

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope
Three walker classes: storefronts (page_size ≤ 50, shallow), partner
syncs (full 100M walks over hours, must resume), batch jobs. Churn
~50 writes/sec means thousands of mutations during one full walk —
stability semantics are the design, not a footnote.

## 2. Wire contract (AIP-132/158)
`GET /v1/products?page_size=100&page_token=...` →
`{ "products": [...], "next_page_token": "CgZwcm9kXzk..." }`
- Default page_size 50, max 1000; server may return fewer than asked
  (a shard was slow) — fewer ≠ done.
- Done is signaled ONLY by empty `next_page_token`.
- page_size 0/absent → default; negative → INVALID_ARGUMENT.

## 3. The token
Opaque base64 of `{ last_sort_key, last_id, query_fingerprint,
minted_at }` — keyset, not offset:
- **Offset dies twice at 100M**: `OFFSET 5,000,000` scans and discards
  5M rows per page (O(n) per page, O(n²) per walk); and a row inserted
  before the cursor shifts every later offset — clients see repeats
  and misses. With 50 writes/sec, a multi-hour walk drifts by ~10^5
  rows.
- **Keyset**: `WHERE (sort_key, id) > (last_sort_key, last_id) ORDER BY
  sort_key, id LIMIT page_size` — O(page) per page via the index; the
  id tiebreaker makes the sort total. No repeats of passed keyspace;
  rows that existed at walk start are all seen; late inserts behind
  the cursor are correctly out of scope.
- Fingerprint binds the token to its filter/order — reuse after
  changing either → INVALID_ARGUMENT ("token does not match request").
- Tokens expire (e.g. 24h via minted_at); garbage → INVALID_ARGUMENT,
  never a silent empty page.

## 4. Evolution & operations
`total_size` omitted on the main List: exact COUNT at 100M per request
is a table scan or a maintained counter nobody needs per-page; an
estimated count lives on a separate stats endpoint if product wants
it. Token opacity is what lets the format evolve (add fields, switch
encoding) with zero client changes (AIP-180). Partner walkers get
higher rate limits but larger mandatory page_size.
