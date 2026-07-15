# My Design: Scaling to Millions of Users

## Step 1 — Use cases, constraints, estimates

<!-- The app, the read/write ratio, and rough requests/sec at each
     user milestone -- these numbers drive every stage. -->

## Step 2 — High-level design (the end state)

<!-- Sketch where the story ends: LB -> autoscaled web tier -> cache ->
     replicated/sharded DB, object store + CDN, queues + workers. -->

## Step 3 — The progression (the heart of this question)

<!-- For each stage: what breaks, and the minimal fix.
     1 box -> split web/DB -> add LB + second web server (why: single
     point of failure and CPU) -> static assets to object store + CDN
     (why: bandwidth) -> DB read replicas (why: read load) -> cache
     layer for hot queries/sessions (why: replica lag/cost) ->
     autoscaling (why: diurnal peaks) -> async queues + workers (why:
     slow requests) -> shard/federate or NoSQL (why: write ceiling).
     Justify each with the numbers from step 1. -->

## Step 4 — Keep it running

<!-- Monitoring/alerting as the bottleneck-finder, backups, multi-AZ
     failover, cost notes. What would 100M users force next? -->
