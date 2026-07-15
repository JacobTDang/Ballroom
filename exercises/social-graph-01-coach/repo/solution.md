# My Design: Social-Network Data Structures

## Step 1 — Use cases, constraints, estimates

<!-- Edges = users x avg friends. Bytes per adjacency list, total
     storage -- show it can't fit on one machine. QPS for lookups vs
     path queries. -->

## Step 2 — High-level design

<!-- Client -> API -> lookup service (user -> person-server) ->
     person servers holding adjacency lists. Where does a path query
     run? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Data model: adjacency lists per user; how are they stored and
       serialized on a person server?
     - Sharding: how do you assign users to person servers? What does
       the lookup service store, and how does it scale?
     - Distributed BFS: the frontier crosses shards every hop -- batch
       the per-server requests, track visited, cap the depth.
     - Bidirectional BFS as the latency fix -- why it helps. -->

## Step 4 — Scale it

<!-- Cache adjacency lists (hot users), replicate person servers for
     reads, handle the heavy tail (a 5M-friend account), and keep the
     lookup service from becoming the bottleneck. -->
