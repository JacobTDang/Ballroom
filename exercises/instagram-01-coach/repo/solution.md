# My Design: Instagram

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope (post, home feed, likes)? What's
     out? Estimate: photos posted/day, storage per photo across
     renditions, feed reads/day, and the read:write ratio -- show the
     arithmetic. -->

## Step 2 — High-level design

<!-- Separate the media path (upload -> object storage + CDN, multiple
     renditions) from the metadata/feed path (post metadata, follow
     graph, feed service). What does posting touch end to end? What
     does opening the feed touch? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Feed generation: fan-out-on-write to followers' feed stores, or
       fan-out-on-read? What happens when a celebrity with millions of
       followers posts?
     - Feed storage: where do materialized feeds live, how long are
       they, and how are they paginated?
     - Counters: how do likes/comments get counted without a
       synchronous UPDATE on every tap? -->

## Step 4 — Scale it

<!-- Where does this break at 10x? How are the feed store and post
     metadata sharded? What does the CDN absorb versus what still hits
     your servers? How big is the fanout worker pool, and from what
     arithmetic? -->
