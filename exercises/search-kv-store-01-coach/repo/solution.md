# My Design: Search Key-Value Cache

## Step 1 — Use cases, constraints, estimates

<!-- Searches/sec from the monthly figure. Entry size x cacheable
     working set -> total memory -> node count. Show the arithmetic. -->

## Step 2 — High-level design

<!-- Client -> LB -> query API -> cache layer -> (miss) search
     backend -> populate cache. Where does the hash ring live? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - LRU at the data-structure level: hash map for O(1) lookup +
       doubly-linked list for O(1) recency updates -- why you need both.
     - Placement: consistent hashing of the query -> cache node; what
       happens when a node joins/leaves.
     - What exactly is the key? Normalized query (case, whitespace,
       term order?) -- cache hit rate depends on this.
     - Freshness: TTL vs explicit invalidation when the index updates. -->

## Step 4 — Scale it

<!-- Hot keys (a breaking-news query hammering one node): local/L1
     caching or key replication. Node failure: what's lost, and why
     that's acceptable for a cache. Memory pressure tuning. -->
