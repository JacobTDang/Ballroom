# My Design: Twitter Timeline & Search

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope? Estimate: tweets/sec written,
     timeline reads/sec, fanout writes per tweet (avg followers),
     storage growth per year. Show the arithmetic. -->

## Step 2 — High-level design

<!-- Separate the write path (post tweet -> fanout -> followers'
     timelines) from the read path (fetch timeline). Where does
     search fit? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Fanout: on write, on read, or hybrid? What happens when a user
       with 100M followers tweets?
     - Timeline storage: where do materialized timelines live, and why?
       How large is one user's timeline allowed to grow?
     - Tweet storage: where does the tweet itself live vs the timeline
       entries that reference it?
     - Search: how are tweets indexed, and what does a query touch? -->

## Step 4 — Scale it

<!-- Where does this break at 10x? Hot users, thundering herds on
     celebrity tweets, cache sizing, sharding the timeline store and
     the search index. -->
