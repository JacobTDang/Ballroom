# Reference Design — Upload 10 GB, Resumably

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
The durable noun is the media resource (`media/{id}`); the working noun
is the **upload session** — a short-lived resource created per upload,
owned by the caller, authorized only for that upload, expiring after
7 days. Mobile drops ~hourly on multi-hour uploads: resume is the
feature, not an edge case. p50 40 MB / max 10 GB drives the tiers.

## 2. Methods & wire shapes
- Tiers: **single-shot** `POST /v1/media?uploadType=media` for < 8 MB;
  **resumable** for everything else (metadata sent at initiation,
  AIP-133-style create).
- Initiate: `POST /v1/media?uploadType=resumable {title, tags}` →
  `200` + session URI (opaque, signed, expiring).
- Chunk: `PUT <session-uri>` with
  `Content-Range: bytes 0-8388607/10737418240` → `308 Resume
  Incomplete` + `Range: bytes=0-8388607` (server-committed bytes).
- Status query after a drop: `PUT <session-uri>` with
  `Content-Range: bytes */10737418240` and an empty body → `308` +
  committed `Range`; resume from the next byte.
- Final chunk → `201` + the media resource JSON.

## 3. The hard part
The **server's committed offset is the only truth** — a client that
lost a response must never guess; it asks (`bytes */N`) and resumes.
Chunks that skip ahead of the committed offset →
`400 INVALID_ARGUMENT` naming the expected offset; overlapping bytes
before it are ignored-by-design only if byte-identical, else the
session fails. **Integrity**: CRC32C per chunk (header) verified on
receipt; whole-file MD5 verified at finalization — any mismatch fails
the session (`FAILED_PRECONDITION`, reason `CHECKSUM_MISMATCH`) and
the client restarts; corrupt bytes never finalize.

## 4. Evolution & operations
Finalization is atomic: the media resource (with metadata from
initiation) exists only after the last byte verifies — no partial
media. Sessions: `DELETE <session-uri>` aborts; abandoned sessions GC
at 7 days. Retry matrix: 5xx/timeout → offset query + chunk retry with
backoff; 4xx invalid-range → client bug, restart; expiry/checksum →
restart session. Chunk size 8-32 MB: at ~1 drop/hour on a 20 Mbit
link, an 8 MB chunk re-uploads ~3 s of work, while sub-1 MB chunks
waste round trips — pick 8 MB mobile / 32 MB server. Thresholds and
chunk sizes are server-advertised so they can evolve without breaking
clients (AIP-180 spirit).
