# Reference Design — YouTube

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Upload + watch + view counts in scope; recommendations, comments, live
streaming, ads out. ~500 hours uploaded/minute → **~720,000 hours/day**
of source video; transcoded into ~5 renditions (adaptive bitrate
ladder) inflates stored bytes several-fold over the source alone. But
**egress dominates**: billions of watches/day, each streaming tens to
hundreds of MB, dwarfs the one-time cost of storing and encoding each
video once — bandwidth, not storage, is the number to design against.

## 2. High-level design
Upload path: client → upload service → raw video to object storage +
a job enqueued → transcoding workers pull the job, produce renditions →
renditions land in object storage, video flips to "live" in the
metadata DB. Watch path: client → metadata lookup → CDN serves the
video segments directly (origin = object storage), never proxied
through app servers.

## 3. Core components
- **Transcoding:** async, queue + worker pool, retries on failure.
  Output is multiple resolution/bitrate renditions split into segments
  (HLS/DASH) so the player can adapt mid-stream to connection quality.
  Video state machine: `uploaded → processing → live` (and `failed`),
  surfaced to the uploader.
- **Serving & CDN:** video popularity is heavily skewed (a small
  fraction of videos is most of the traffic) — exactly the shape a CDN
  is for. Long tail cold videos fall back to origin object storage
  on a cache miss. Thumbnails are a separate small-object problem
  (many tiny files, different caching characteristics than video
  segments).
- **View counts:** incremented via a buffered/batched aggregator, not
  a synchronous UPDATE per view — a viral video's view row would
  otherwise become a massive write hotspot. Displayed count is
  eventually consistent with actual views.

## 4. Scale
- Transcoding fleet auto-scales with upload queue depth, independent of
  watch traffic entirely.
- CDN absorbs almost all read/watch scale; origin storage mainly serves
  cache misses on cold/old videos.
- Metadata DB sharded (by video ID), cached for hot videos; cold,
  rarely-watched video renditions tiered to cheaper storage classes
  over time instead of kept on hot storage forever.
