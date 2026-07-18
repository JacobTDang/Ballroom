# Design the Library API

Design the public HTTP API for a lending library: publishers, their
books, physical copies on shelves, and members who check copies out
and return them. This is the resource-modeling question — every later
question in the track builds on the discipline you set here.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Which nouns are resources? Publisher, book, copy, member, loan —
  which of these deserve their own collection, and which are fields
  on something else?
- Resource names: flat (`books/{book}`) or hierarchical
  (`publishers/{publisher}/books/{book}`)? What does the hierarchy
  buy you, and what does it cost when a book changes publisher?
- Who calls this API? Branch software, a public catalog site, partner
  integrations — does that change the surface?

## Suggested defaults (if you want a starting point)

- One library system, many branches; ~1M books, ~5M copies,
  ~500k members.
- Checkout and return must be first-class operations — the librarian
  desk hits them constantly.
- Search/browse is read-heavy; lending writes are modest.

## What good looks like

By the end you should have: the resource hierarchy with concrete
names, the five standard methods (Get/List/Create/Update/Delete)
mapped to HTTP for at least books and copies, a deliberate decision on
how checkout/return work (custom `:checkout`-style methods or a loan
resource — argued, not assumed), request/response shapes for the core
calls with real field names, list pagination at least sketched
(`page_size`/`page_token`), and the obvious error cases (404 vs 400
vs 409) handled consistently.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
