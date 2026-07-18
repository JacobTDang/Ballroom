# Grading Rubric — Upload 10 GB, Resumably

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Protocol tiering

- Single-shot vs resumable (and optionally multipart-with-metadata)
  distinguished, with concrete size/reliability thresholds and a
  reason for each.
- Small files are not taxed with session overhead.

## 2. Session model

- Initiation returns a session handle/URI; the session is a
  short-lived resource with an owner, an expiry, and auth scoped to
  that one upload.
- Session state is inspectable (status query), not write-only.

## 3. Offset & resume semantics

- The server is authoritative about committed bytes; the client's
  recovery step after a failure is exactly "query committed offset,
  resume from there" (the 308 + Range dance or an equivalent).
- Out-of-order and overlapping chunks are rejected loudly, not
  silently reassembled.

## 4. Integrity

- A checksum story: per-chunk and/or whole-file (CRC32C/MD5-class),
  stated where it is computed and verified.
- A mismatch fails the session loudly; the design says what the
  client does next (restart), never silently accepts corrupt bytes.

## 5. Finalization & lifecycle

- A definite finalization moment produces the durable media resource
  (with its metadata); partial sessions never appear as media.
- Abandoned sessions expire and are GC'd on a stated window; explicit
  abort exists.

## 6. Failure & retry matrix

- Failures classified: transient network/5xx retries the chunk (after
  an offset query); session-fatal errors (expiry, checksum mismatch,
  invalid range) restart the session.
- Error shapes are consistent with a canonical error model (AIP-193),
  not ad-hoc strings.

## 7. Communication & trade-offs

- Chunk size argued from the numbers (throughput vs re-upload cost on
  a drop); resumable overhead for small files acknowledged.
- Trade-offs stated for the major choices; the candidate drove the
  4 steps themselves.
