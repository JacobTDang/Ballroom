# My Design: URL Shortener

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope? What are you assuming about
     traffic, read/write ratio, and link lifetime?
     Estimate: writes/sec, reads/sec, storage for 3 years, bandwidth.
     Show the arithmetic — the numbers drive every later choice. -->

## Step 2 — High-level design

<!-- Boxes and arrows, end to end: what does a shorten request touch?
     What does a redirect touch? Keep it to the major components. -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Short-code generation: how is it derived? Why won't two URLs
       collide? How long does the code need to be for your scale?
     - API design: the endpoints, their inputs/outputs.
     - Data model: tables/collections, and SQL vs NoSQL — justify it
       with your access patterns from step 1.
     - The redirect: which HTTP status code, and why? -->

## Step 4 — Scale it

<!-- Where does this break at 10x the load from step 1? Which reads
     should be cached and with what policy? How does the database
     scale — replication, sharding? What about hot links? -->
