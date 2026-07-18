# Grading Rubric — Dropbox

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped to upload/download, cross-device sync, and sharing; explicit
  out-of-scope calls (collaborative editing, previews).
- Estimates: users × storage each, upload bandwidth, and the read/write
  shape; dedup savings called out as a storage lever.

## 2. High-level design

- Metadata service (file tree, versions, block lists) cleanly separated
  from block storage (the bytes, in object storage); clients talk to
  both through an API layer.

## 3. Chunking & dedup

- Files split into content-addressed blocks (hash as ID); only changed
  blocks upload on edit (delta sync); identical blocks deduped across
  users, with the privacy caveat acknowledged.

## 4. Sync protocol

- Client watcher detects local changes; server notifies other devices
  (long-poll or push) rather than devices polling blindly; the full
  round trip narrated: edit on laptop → blocks + metadata up → notify →
  phone pulls metadata → fetches missing blocks.

## 5. Conflicts & consistency

- Concurrent edits handled deliberately: version vectors or
  last-writer-wins plus a conflicted-copy, with the choice defended.
- Metadata is the source of truth; a block upload without its metadata
  commit is harmless.

## 6. Scaling story

- Metadata DB sharded (by user/namespace); hot blocks cached; cold
  blocks tiered to cheaper storage; notification fanout scaled
  separately from data transfer.

## 7. Communication & trade-offs

- Trade-offs stated (block size vs overhead, sync freshness vs
  battery/bandwidth), driven by the estimates; the candidate drove the
  4-step structure.
