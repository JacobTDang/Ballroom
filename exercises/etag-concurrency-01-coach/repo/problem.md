# Two Admins, One Settings Page

An organization's settings live behind your API and are edited from a
web console. Two admins open the same settings page, both edit, both
save. Today the last writer silently wins — the first admin's change
vanishes with no error, no trace. Design the read/update surface so
that can't happen.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go. The canon is optimistic concurrency
with etags (AIP-154, HTTP If-Match/412) — derive the mechanism from
the failure mode before naming it.

## Scope to establish in step 1

Talk through and pin down with your coach:

- The resource: a settings singleton per org? How big, how often
  written, how often read?
- Who edits: humans in a console (seconds-to-minutes between read
  and write) and automation (read-modify-write in milliseconds)?
- What should the losing admin experience — an error? a merge? and
  who does the merging?

## Suggested defaults (if you want a starting point)

- One settings document per org, ~5 KB, dozens of writes/day
- Read on every console page load; a few automations also write
- A stale save should never silently destroy the other admin's edit

## What good looks like

By the end you should have: the lost-update interleaving written out
concretely; etags on reads and required (or strongly recommended)
If-Match on writes with 412 on mismatch — and a stated choice of what
the etag is derived from; the strong-vs-weak validator distinction and
which this API needs; the read-side If-None-Match/304 win; a concrete
client recovery protocol after a 412; and a policy call on whether
unconditional writes are allowed at all.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
