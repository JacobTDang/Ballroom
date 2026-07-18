# My Design: Nearby Friends

## Step 1 — Use cases, constraints, estimates

<!-- Who sees whom, and under what opt-in condition? What's out of
     scope? Estimate: active opted-in users x ping interval ->
     location writes/sec (this should be the defining, write-heavy
     number), and the state one user's location needs. -->

## Step 2 — High-level design

<!-- Mobile clients ping location periodically -> location service ->
     in-memory location store + geo index. What's the query/notify path
     that actually tells a user which friends are near? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Geo-indexing: geohash buckets or a quadtree? What's the
       cell-size trade-off, and how do you handle two people who are
       close in reality but land in adjacent cells?
     - Push vs pull: pub/sub pushing updates to nearby friends'
       devices, or periodic pull? Work through the friend-pair fanout
       cost either way.
     - Privacy: where exactly does per-friend-pair opt-in get checked? -->

## Step 4 — Scale it

<!-- Where does this break at 10x -- especially in a dense urban area?
     How is the location store sharded? What ages out a stale/offline
     user so they don't show as a ghost? Can update frequency adapt to
     save battery and write load? -->
