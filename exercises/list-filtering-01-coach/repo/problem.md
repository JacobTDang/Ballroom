# Filter and Order the Fleet

A device-management platform tracks two million IoT devices, and every
customer dashboard starts with the same call: List devices, filtered
and sorted their way. Design the filtering and ordering surface for
that List endpoint — expressive enough to be useful, bounded enough
that your datastore can actually serve every query you promise.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- What do operators actually filter by? Status, firmware version,
  site, last-seen time, tags — which fields, which operators?
- What must sorting support — single key or multi-key? What's the
  default order when the client says nothing?
- Is free-text search ("find the device named like...") the same
  feature or a different one?

## Suggested defaults (if you want a starting point)

- 2M devices across 10k customers; the largest single customer has
  200k devices.
- Filterable: status (enum), site, firmware_version, last_seen_time,
  labels. Sortable: name, last_seen_time, firmware_version.
- The datastore is indexed, not magic — arbitrary filter × arbitrary
  sort is not free.

## What good looks like

By the end you should have: a bounded filter grammar (AIP-160 style —
fields, operators, AND at minimum) with unfilterable fields rejected
loudly, an `order_by` contract (whitelisted keys, directions, stable
default), an honest mapping from promised queries to indexes with the
combinatorial explosion named and bounded, the pagination interaction
pinned (tokens bound to filter+order), a clear split between exact
filtering and relevance-ranked search, and malformed-filter errors
that say what's wrong and where.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
