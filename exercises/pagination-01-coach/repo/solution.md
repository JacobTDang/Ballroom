# My Design: Catalog Pagination

## Step 1 — Scope & resource model

<!-- Who walks the catalog and how (interactive vs full sync)? How
     much churn happens underneath them? Which orderings must exist?
     Put numbers on it: rows, churn/sec, walk duration. -->

## Step 2 — Methods & wire shapes

<!-- The List call per AIP-132/158: request (page_size, page_token,
     order_by?) and response (products, next_page_token). Defaults,
     maximum, server-may-return-fewer, zero/absent handling. One
     worked example of page 1 → page 2. -->

## Step 3 — The hard part: the token

<!-- What exactly is inside next_page_token, and why is it opaque?
     Keyset position vs offset — demonstrate with numbers why OFFSET
     dies at 100M. What happens to a walker when rows are inserted or
     deleted behind/ahead of it? Expiry, tampering, and a token
     reused after the client changed filters. -->

## Step 4 — Evolution & operations

<!-- total_size: offer it or not, and the counting cost. What can
     change in the token format later (it's opaque — prove that
     matters)? Rate limits for full-walk partners vs storefronts. -->
