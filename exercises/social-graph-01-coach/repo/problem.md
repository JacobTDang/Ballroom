# Design Data Structures for a Social Network

Design the data layer for a social network's friend graph: store who
is connected to whom, and answer "show me the shortest chain of
friends between me and this person" — at a scale where the graph does
not fit on one machine.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: look up a person's friends, find the shortest path
  between two people (degrees of separation). Out of scope? (Feed,
  posts, suggestions.)
- The defining constraint: the graph is too big for one box. Every
  design decision follows from that.
- Put numbers on it: users, average friends per user, edges total,
  bytes per adjacency list, machines needed.

## Suggested defaults (if you want a starting point)

- 100 million users, ~1,000 friends each (heavy-tailed)
- ~100 billion friend edges total
- Shortest-path queries are common but latency-tolerant (~1s is fine);
  friend lookups must be fast
- The graph is mostly-read; edge writes (friend/unfriend) are rare by
  comparison

## What good looks like

By the end you should have: sizing arithmetic that proves the graph
must be partitioned, a person-server sharding scheme with a lookup
service mapping user → shard, an adjacency-list data model, a BFS that
works when every hop may live on a different machine (and what that
does to latency), and a scaling story covering caching, replication,
and the heavy-tailed degree distribution.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
