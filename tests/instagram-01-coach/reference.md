# Reference Design — Instagram

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Post photo + home feed + likes in scope; stories, reels, DMs, explore
out. 100M photos/day → **~1,000 writes/s** (86,400 s/day rounds to
100K). Feed reads ~1000:1 over posts → **~1M reads/s**. Each photo
transcoded into ~4 renditions (thumbnail/feed/full/original);
read-heavy by a wide margin, so the feed read path is what step 2
onward optimizes for.

## 2. High-level design
Upload: client → media service → object storage (original) → async
transcode workers → renditions on object storage + CDN. Metadata: post
service writes post row + triggers fanout. Read: client → feed service
→ materialized feed store, media served straight from the CDN (never
proxied through app servers).

## 3. Core components
- **Fanout-on-write** for normal accounts: on post, push the post ID
  into every follower's feed list (capped length, e.g. ~800 entries).
- **Celebrity problem:** don't fan out for accounts above a follower
  threshold; merge their recent posts in at read time instead. This
  hybrid is the expected answer — pure fan-out-on-write means one post
  triggers millions of writes; pure fan-out-on-read means every normal
  feed load re-merges hundreds of following-lists.
- **Feed storage:** materialized feeds as ID lists in a fast store
  (e.g. Redis), cursor-paginated; cold/inactive users' feeds rebuilt
  on-demand from the follow graph instead of kept warm forever.
- **Counters:** likes/comments incremented via a buffered/async
  aggregator (sharded counters or a queue that periodically flushes),
  not a synchronous UPDATE per like — accept brief eventual consistency
  on the displayed count.

## 4. Scale
- Feed store sharded by user, post metadata sharded by post ID; CDN
  absorbs the read-heavy media traffic entirely.
- Fanout worker pool sized from posts/day × avg-follower-count, running
  behind a queue so a burst of posts doesn't block the write path.
- Hot posts (viral content) cached at the metadata layer too, not just
  media, since a fan-out-on-read merge would otherwise re-fetch them
  constantly.
