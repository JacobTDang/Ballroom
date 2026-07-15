# Grading Rubric — Social-Network Data Structures

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to friend lookups + shortest path, with explicit out-of-scope
  calls (feed, posts, recommendations).
- Stated the defining constraint unprompted: the graph exceeds one
  machine, so partitioning drives everything.

## 2. Back-of-envelope estimates

- Edge count from users × average degree; adjacency-list bytes and
  total storage shown to exceed a single box.
- Latency budgets distinguished: fast friend lookups vs slower path
  queries.

## 3. High-level design

- Person servers holding adjacency lists, fronted by a lookup service
  mapping user → shard; clean client → API → graph-layer path.

## 4. Data model & sharding

- Adjacency lists as the core structure, keyed by user.
- A concrete shard assignment (hash(user) or range) and where the
  mapping lives; how a shard split/rebalance would work.

## 5. Distributed shortest path

- BFS that batches frontier lookups per person server each hop, with a
  shared visited set and a depth cap.
- Bidirectional BFS (or an argued equivalent) named as the mitigation
  for cross-shard hop cost.

## 6. Consistency & writes

- Friend/unfriend writes update both users' adjacency lists; the
  design says how the two-sided update stays consistent (transaction,
  or async with reconciliation).

## 7. Scaling story

- Hot adjacency lists cached (LRU) and person servers replicated for
  reads; the heavy tail (multi-million-degree accounts) handled
  explicitly (chunked lists, pagination).
- Lookup service kept cheap: cacheable, replicated, tiny records.

## 8. Communication & trade-offs

- Trade-offs stated (precompute vs on-demand paths, memory vs disk for
  adjacency lists), driven by the step-1 numbers; the candidate drove
  the 4-step structure.
