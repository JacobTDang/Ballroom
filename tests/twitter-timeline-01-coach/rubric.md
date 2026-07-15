# Grading Rubric — Twitter Timeline & Search

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to timeline + posting + search, with explicit out-of-scope
  calls (DMs, notifications, trending).
- Identified the celebrity/fanout asymmetry as the defining constraint
  without being told.

## 2. Back-of-envelope estimates

- Tweets/sec and timeline-reads/sec derived from assumed volume, with
  the read-heavy ratio stated and used later.
- Fanout arithmetic: average followers × tweets/sec = timeline writes/sec,
  and what that number becomes for a 100M-follower account.

## 3. High-level design

- Write path (tweet service → fanout service → timeline stores) drawn
  separately from the read path (timeline service → cache).
- Search cluster fed from the write path; nothing load-bearing appears
  later that isn't on the diagram.

## 4. Fanout strategy

- A concrete decision: fan-out-on-write for normal users, with a
  deliberate answer for celebrities (fan-out-on-read for high-follower
  accounts, merged at read time — or an argued alternative).
- Consequences stated: write amplification vs read latency trade-off.

## 5. Timeline & tweet storage

- Materialized timelines in memory (e.g. Redis lists) with a bounded
  length per user, justified by the latency target.
- Tweets stored once (SQL/NoSQL + object store for media); timelines
  hold references, not copies.

## 6. Search

- Tweets tokenized and written to an inverted index (or an argued
  service like Elasticsearch); query path described end to end.
- Freshness expectation stated (how soon is a new tweet searchable?).

## 7. Scaling story

- Timeline store sharded by user; search index sharded by term or
  document, with the choice justified.
- Cache-first reads; thundering-herd handling for viral tweets.

## 8. Communication & trade-offs

- Trade-offs stated for fanout, storage, and index choices, driven by
  the step-1 numbers; the candidate drove the 4-step structure.
