# My Design: Web Crawler

## Step 1 — Use cases, constraints, estimates

<!-- Pages/month -> fetches/sec (show the arithmetic). Storage for
     extracted text at your page size. Refresh cadence. -->

## Step 2 — High-level design

<!-- The crawl loop as components: URL frontier (queue), fetchers,
     parser/extractor, dedup, storage, back into the frontier. -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - URL dedup: how do you know you've seen a URL before, at
       billions of URLs? Exact set, hashes, bloom filter trade-offs.
     - Content dedup: two different URLs, same page -- how do you
       detect it (signature/simhash) and what do you do?
     - Politeness: how does the frontier guarantee one fetcher at a
       time per domain, with rate limits and robots.txt?
     - Crawler traps: infinite calendars, session-id URLs -- defenses? -->

## Step 4 — Scale it

<!-- Distribute the frontier: partition by domain (why?), consistent
     hashing across crawler nodes, prioritization/freshness queues,
     recovering from a dead crawler node. -->
