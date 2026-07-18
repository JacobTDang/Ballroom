# Bulk Import for Contacts

A CRM exposes a contacts collection with the standard per-contact
methods (get, list, create, update, delete). Customers now sync from
spreadsheets and other CRMs: thousands of contacts per run, and one
HTTP call per contact is too slow and too chatty. Design the batch
surface.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go. The canon is the batch-method AIPs
(231/233/234/235): batch methods mirror their singular methods — try
to derive why that mirroring matters before leaning on it.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Which batch methods earn their keep: batchGet? batchCreate?
  batchUpdate? batchDelete? All four, or fewer?
- The resource model stays the same — what changes is the request
  envelope. What does one batch item contain?
- Typical import size, payload budget, and how often imports repeat
  (the same spreadsheet, re-run).

## Suggested defaults (if you want a starting point)

- Typical import: 10,000 contacts; largest: 100,000
- Max request payload: 10 MB
- Imports get re-run after partial failures (assume retries happen)

## What good looks like

By the end you should have: batch methods whose requests embed the
singular requests (not a bespoke import format); an explicit atomicity
decision per method — all-or-nothing or partial success — with the
reasoning; per-item results position-matched to the request with no
silent drops; stated limits (items and bytes) and the over-limit
error; a story for retrying a half-succeeded batch without duplicating
contacts; and the honest arithmetic for when a batch API stops being
enough and an import job takes over.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
