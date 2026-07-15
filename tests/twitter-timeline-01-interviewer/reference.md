# Reference Design — Twitter Timeline & Search

## 1. Use cases, constraints, estimates
Post tweet, home timeline, user timeline, search; DMs/ads out. 100M
DAU, 500M tweets/day → **~6K writes/s** (peaks 10×). Fanout: ×200 avg
followers → **~1.2M timeline writes/s** — fanout, not tweet storage, is
the hard number. Timeline reads dominate; target < 200 ms.

## 2. High-level design
Write path: client → LB → tweet write API → tweet store, then **fanout
service** → per-user timeline caches + search ingester.
Read path: client → LB → timeline API → timeline cache (fast path).

## 3. Core components
- **Fanout on write** for normal users: push tweet_id into each
  follower's timeline list (Redis, capped ~800 entries/user).
- **Celebrity problem:** don't fan out for high-follower accounts; pull
  their recent tweets at read time and merge — the hybrid is the
  expected answer.
- **Storage:** tweets once in a sharded store (id, user, text, media
  ref); timelines hold IDs only; media in object storage + CDN.
- **Search:** tweets tokenized on write into an inverted index
  (sharded, e.g. Elasticsearch-like); freshness within seconds.

## 4. Scale
- Timeline cache sharded by user; tweet store by tweet id; search index
  by term or doc — state which and why.
- Thundering herd on viral tweets: cache the tweet object itself
  everywhere, coalesce requests.
- Fanout worker pool sized from writes/s × avg followers; queue between
  write API and fanout absorbs spikes.
