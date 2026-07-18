# Reference Design — Dropbox

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Upload/download, cross-device sync, sharing (link/invite); real-time
collaborative editing and in-browser previews out. 500M users × ~2 GB
avg → **~1 EB** stored. ~10% DAU sync something → 50M syncers/day, a
few files each — **modest write QPS**, the story here is storage and
bandwidth, not request rate. Dedup across users on identical files
(stock installers, common docs) is a real lever, not a footnote.

## 2. High-level design
Client → API layer → two backends: **metadata service** (file tree,
versions, per-file block list) on a DB, and **block storage** (the
bytes) in object storage, addressed by content hash. Clients diff
locally, upload changed blocks + a metadata commit.

## 3. Core components
- **Chunking:** fixed or content-defined ~4 MB blocks; block ID =
  hash(content). Only blocks that changed re-upload on an edit (delta
  sync) — a 1-byte change to a 1 GB file uploads ~4 MB, not 1 GB.
- **Dedup:** identical blocks (same hash) across any user's files are
  stored once, reference-counted; note the privacy trade-off (server
  can tell two users have the same bytes) and that it's opt-out-able.
- **Sync protocol:** client watcher detects local changes → uploads
  blocks + commits new metadata → server notifies other online devices
  (long-poll or push) → they pull the new metadata → fetch only the
  blocks they're missing. Metadata commit is the atomic "this version
  now exists" moment; a block upload alone is inert.
- **Conflicts:** metadata versioned (version vector or a simple
  last-writer-wins on the metadata row); a losing concurrent edit is
  never dropped silently — saved as a `filename (conflicted copy).ext`
  for the user to reconcile.

## 4. Scale
- Metadata DB sharded by user/namespace — most access is one user's
  own tree, so this keeps queries local to a shard.
- Hot/recently-accessed blocks cached at the edge; cold blocks tiered
  to cheaper storage classes over time.
- Notification fanout (telling N devices "something changed") scaled
  as its own service, independent from the bulk block-transfer path —
  a notification is small, the payload behind it is not.
