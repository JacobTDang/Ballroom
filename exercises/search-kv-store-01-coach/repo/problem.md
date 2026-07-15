# Design a Key-Value Cache for a Search Engine

Design the cache that sits in front of a search engine's query
pipeline: given a query, return recently computed results without
re-running the expensive search — evicting the right things when
memory runs out.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: look up cached results for a query, store fresh
  results after a cache miss, evict when full, expire stale entries.
  Out of scope? (The search/ranking pipeline itself.)
- The two sizing questions that shape everything: does the working set
  fit in memory, and on how many machines?
- Put numbers on it: searches/sec, hit-rate assumption, bytes per
  cached entry, total cache memory, nodes needed.

## Suggested defaults (if you want a starting point)

- 10 billion searches per month; popular queries are a tiny fraction
  of distinct queries (classic Zipf)
- A cached entry is the query string plus a serialized results page,
  ~250 KB... or much smaller if you argue for caching doc IDs only
- Results may be stale after ~1 hour (index updates)
- Latency target for a hit: a few milliseconds

## What good looks like

By the end you should have: memory arithmetic that determines the node
count, an LRU design you can explain at the data-structure level (hash
map + doubly-linked list and why both), a consistent-hashing scheme for
spreading queries across cache nodes, an expiry/invalidation answer for
index freshness, and a scaling story for hot keys and node failures.
