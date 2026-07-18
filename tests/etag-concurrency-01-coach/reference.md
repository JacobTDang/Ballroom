# Reference Design — Two Admins, One Settings Page

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
One settings document per org: `orgs/{org}/settings`, ~5 KB, dozens of
writes/day, read on every page load. The failure: A GETs v7, B GETs
v7, B PATCHes (v8), A PATCHes — A's write is built on v7 and silently
erases B's v8. Human edit gaps are minutes, so the race window is
wide open; this is the textbook lost update.

## 2. Methods & wire shapes
- `GET /v1/orgs/{org}/settings` → body includes `"etag": "\"v8\""`,
  also sent as the `ETag` header.
- `PATCH /v1/orgs/{org}/settings` with header `If-Match: "v8"` (or
  the body's etag field, AIP-154 — support both, header canonical).
  - Match → 200, new etag `"v9"`.
  - Stale → **412 Precondition Failed** / FAILED_PRECONDITION with
    ErrorInfo reason `ETAG_MISMATCH`.
  - Missing If-Match → 428-style precondition-required error: this
    resource refuses unconditional writes.

## 3. The hard part
- **ETag source**: a server-side monotonic version counter — cheap,
  strong, and never collides; a content hash is the alternative
  (idempotent no-op writes keep the same etag) at hashing cost.
  Either way the value is opaque; clients never parse it.
- **Strong vs weak**: strong validators guarantee byte-identical
  representations — required here, because a write decision hangs on
  it. Weak (`W/"..."`) tolerates semantically-equivalent variation
  and belongs to cache freshness, not write preconditions.
- **Read side**: console sends `If-None-Match: "v8"` on reload →
  `304` empty body when unchanged — at a page load per edit session,
  most reads become 304s and the 5 KB payload disappears.
- **Recovery (documented contract)**: on 412 → re-GET (fresh etag +
  current doc) → auto-rebase non-conflicting field edits and retry
  once; a genuine field-level conflict surfaces both values to the
  human. Automations retry the read-modify-write loop with backoff,
  capped.

## 4. Evolution & operations
Unconditional writes stay forbidden for settings (small write rate,
huge blast radius); a bulk migration tool gets a scoped admin
override, logged. Pessimistic locks rejected: nothing holds a lock
across a human's coffee break. Where etags are overkill: append-only
or single-writer resources (audit events, one importer) — concurrency
control there is ceremony without a race. List responses carry
per-item etags so a console listing many orgs can edit any of them
conditionally (AIP-154's list note).
