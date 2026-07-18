# Grading Rubric — YouTube

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases, constraints & estimates

- Scoped to upload + watch (+ view counts); explicit out-of-scope
  calls (recommendations, comments, live streaming, ads).
- Estimates: hours uploaded per minute → storage/day across encodings;
  the observation that egress bandwidth, not storage, dominates cost.

## 2. High-level design

- The upload pipeline (upload service → queue → transcoding workers →
  object storage) drawn separately from the watch path (CDN → player),
  with metadata DB serving both.

## 3. Transcoding pipeline

- Videos transcoded asynchronously into multiple resolutions/bitrates
  (adaptive streaming, e.g. HLS/DASH segments); queue + worker pool
  with retries; video state machine (uploaded → processing → live).

## 4. Serving & CDN

- Video segments served from a CDN with origin in object storage; the
  popularity skew argument (a tiny fraction of videos is most of the
  traffic) used to justify it; thumbnails handled as their own small-
  object problem.

## 5. Metadata & view counts

- Metadata (title, channel, segment manifest) in a sharded DB, cached
  for hot videos.
- View counts aggregated asynchronously (buffered/batched increments),
  not synchronous row updates on every view — the write-hotspot
  reasoning stated.

## 6. Scaling story

- Transcoding fleet scales with upload volume, CDN absorbs read scale;
  hot-video cache; storage tiering for old/cold videos.

## 7. Communication & trade-offs

- Trade-offs stated (encoding cost vs playback quality ladder,
  count accuracy vs write load), driven by the estimates; the candidate
  drove the 4-step structure.
