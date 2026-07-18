# Paginate the Product Catalog

An e-commerce platform's catalog holds 100 million products, and every
storefront, feed, and partner sync walks it through your List API.
Design the pagination: the token format, the guarantees, and the
failure modes. This question is where "just add OFFSET" goes to die.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Who paginates? Interactive storefronts (small pages, shallow),
  partner syncs (full walks, resumable), internal batch jobs.
- How live is the data under the walkers? Products are created,
  deleted, and re-priced constantly.
- What ordering do clients actually need — and what does each
  ordering cost you?

## Suggested defaults (if you want a starting point)

- 100M products, ~50 writes/sec of churn.
- A full partner walk takes hours and must be resumable.
- Default ordering: stable by product id; newest-first exists as a
  secondary need.

## What good looks like

By the end you should have: the `page_size`/`page_token` contract per
AIP-158 (defaults, max, server-may-return-fewer), an opaque
server-minted token whose contents you specify (keyset position +
what else?), the concrete argument with numbers for why OFFSET fails
at this scale, pinned behavior when rows are inserted/deleted
mid-walk, token lifecycle (expiry, tamper handling, what happens if
the client changes filters mid-token), and a deliberate call on
`total_size`.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
