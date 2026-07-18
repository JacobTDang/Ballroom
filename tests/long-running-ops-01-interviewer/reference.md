# Reference Design — Transcode Video: Long-Running Operations

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
Videos already exist (`videos/{id}`). Transcoding is asynchronous work,
so it gets its own resource: an **operation** (`operations/{id}`),
created by the start call and owned by the requesting project. 50k
jobs/day ≈ 0.6 starts/sec — the load is all in the polling, not the
starting. Console users watch one job; pipelines track thousands.

## 2. Methods & wire shapes
- Start: `POST /v1/videos/{id}:transcode {profile: "1080p-h264"}` →
  returns the operation immediately (AIP-151):
  `{"name":"operations/op123","done":false,"metadata":{...}}`
- Poll: `GET /v1/operations/op123`.
- List: `GET /v1/operations?filter=video=videos/v42` (AIP-160-lite).
- Finished success:
  `{"name":"...","done":true,"response":{"renditions":[...]}}`
- Finished failure:
  `{"name":"...","done":true,"error":{"code":13,"message":"..."}}`
  — response XOR error, never both (AIP-151).

## 3. The hard part
- **Metadata vs response**: metadata is the in-flight channel —
  `{progress_percent: 40, stage: "AUDIO_MUX", eta: "..."}` — advisory,
  updated freely. The response appears only at `done: true` and never
  changes afterward.
- **Terminal states**: success (response), failure (error with a
  canonical code, AIP-193), cancelled (error code CANCELLED). Clients
  branch on `done` then on which result field is set.
- **Cancel**: `POST /v1/operations/op123:cancel` — best-effort: the
  server flags the job; a job that finishes anyway stays successful.
  Cancel is idempotent; cancelling a done operation is a no-op
  returning the operation (AIP-151's stance), documented.
- **Notification**: polling contract — start at 1s, exponential
  backoff to a 30s cap; a poll is a cheap point read. Push (webhook or
  queue) halves latency and load but brings delivery/retry/endpoint
  management — deferred to the eventing surface (its own design), the
  right call at 50k/day.

## 4. Evolution & operations
Finished operations retained **30 days**, then GET → NOT_FOUND
(documented; the renditions live on the video regardless — losing the
operation id costs history, not results). Metadata is additive-only:
new fields (e.g. per-rendition progress) never repurpose old ones
(AIP-180). At 10x volume the poll fleet grows linearly — the push
option becomes worth its complexity around there; the operation shape
does not change, which is the point of the pattern.
