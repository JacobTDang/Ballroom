# Reference Design — Bulk Import for Contacts

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
Contacts stay exactly as they are; batching changes the envelope, not
the resource. Ship `batchCreate`, `batchUpdate`, and `batchGet`
(imports and syncs need them); `batchDelete` optional. 10k-contact
imports at ~1 KB/contact ≈ 10 MB — exactly the payload budget, so
limits land well below it.

## 2. Methods & wire shapes
- `POST /v1/contacts:batchCreate`
  `{"requests": [{"contact": {...}}, {"contact": {...}}]}` →
  `{"contacts": [{...}, {...}]}` — each embedded request is exactly a
  CreateContact request (AIP-233); same for `:batchUpdate` (AIP-234,
  each item carrying its update mask).
- `GET /v1/contacts:batchGet?names=contacts/a&names=contacts/b` →
  atomic: all found or the call fails (AIP-231).

## 3. The hard part
- **Atomicity**: batchGet atomic by definition. Mutations:
  **all-or-nothing** — an import wants "it worked" or "fix and
  re-run", and atomic batches make retries trivially safe. (Partial
  success is the defensible alternative; it must then return
  position-matched per-item results, each slot a contact or a
  google.rpc.Status — argued, not defaulted into.)
- **Errors**: atomic failure returns one canonical error whose
  details name the failing indices (AIP-193 BadRequest field
  violations with `requests[17].contact.email`); nothing silent.
- **Limits**: 200 items or 4 MB per batch, whichever first —
  over-limit is INVALID_ARGUMENT naming the limit; clients chunk
  10k-contact imports into 50 batches.
- **Retry dedup**: client-assigned `contact_id` on create (AIP-133)
  makes re-running an import idempotent — replayed creates of the
  same id return ALREADY_EXISTS per item (or succeed-as-noop,
  chosen and documented); alternatively a batch-level request_id
  (AIP-155) replays the stored outcome.

## 4. Evolution & operations
Throughput: 10k contacts = 50 batched calls vs 10,000 singular —
two orders of magnitude fewer round trips; at ~5 batches/sec an
import lands in ~10 s. Ceiling: 100k+ contacts or multi-minute
processing → `contacts:import` returning a long-running operation
(the LRO pattern) with per-item error files, instead of holding a
giant synchronous request. Because batch items ARE the singular
requests, every field added to Contact or CreateContact flows into
batching for free (AIP-180's additive-only rules apply once, not
twice).
