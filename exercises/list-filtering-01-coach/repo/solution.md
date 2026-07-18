# My Design: Fleet Filtering & Ordering

## Step 1 — Scope & resource model

<!-- Who queries, and for what? Which fields must filter, which must
     sort, and what does the biggest customer's 200k-device list mean
     for your worst case? -->

## Step 2 — Methods & wire shapes

<!-- The List request: filter and order_by parameters (AIP-160/132)
     alongside page_size/page_token. Write the grammar you accept —
     fields, operators, conjunction — and three example filter
     strings from real dashboard needs. -->

## Step 3 — The hard part: promises vs indexes

<!-- Map every promised filter/sort combination to an index. Name the
     explosion (F filters × S sorts) and how you bound it. What
     happens to the page token when filter or order changes? Where
     does full-text search live, and why is it not this parameter? -->

## Step 4 — Evolution & operations

<!-- Adding a filterable field later: what has to happen before you
     may promise it? Malformed filters: the exact error shape with
     position/reason. Cost limits on expensive queries. -->
