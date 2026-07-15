# My Design: Mint.com

## Step 1 — Use cases, constraints, estimates

<!-- Transactions/month across all users -> writes/sec (show the
     arithmetic -- note how modest it is). Storage per transaction and
     over 5 years. Read pattern: how often do users actually look? -->

## Step 2 — High-level design

<!-- Separate the ingestion pipeline (daily pulls from financial
     institutions -> categorize -> store -> budget check) from the
     read path (user opens app -> monthly overview). -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Account sync: how do daily pulls work at 30M accounts? Queue of
       sync jobs, workers, handling slow/failed institutions.
     - Categorization: seller -> category rules, what happens for
       unknown sellers, user overrides that stick.
     - Budgets & alerts: where does "you crossed your grocery budget"
       get computed -- on write, on read, or async? Notification path.
     - Data model: transactions table growth, monthly rollups. -->

## Step 4 — Scale it

<!-- Precompute monthly spending per category (batch/MapReduce style)
     instead of scanning transactions on read. Shard the transaction
     store. Security/privacy notes for bank credentials. -->
