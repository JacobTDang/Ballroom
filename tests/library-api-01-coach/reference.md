# Reference Design — Library API

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Scope & resource model
Resources: `publishers/{publisher}`, `publishers/{publisher}/books/{book}`
(the catalog entry, owned by its publisher — AIP-121/122),
`copies/{copy}` (a physical item; flat, because copies move between
branches and outlive shelf assignments), `members/{member}`, and
`loans/{loan}` as the record of a checkout. Shelf is a field on copy,
not a resource — nothing operates on shelves independently.

## 2. Methods & wire shapes
Standard methods per AIP-131..135, e.g. books:
- `GET /v1/publishers/{publisher}/books/{book}` — Get (idempotent)
- `GET /v1/publishers/{publisher}/books` — List, with `page_size`,
  `page_token` → `{ "books": [...], "next_page_token": "..." }`
- `POST /v1/publishers/{publisher}/books` — Create; body is the book;
  response is the created book including its server-assigned name
- `PATCH /v1/publishers/{publisher}/books/{book}` — Update (partial)
- `DELETE /v1/publishers/{publisher}/books/{book}` — Delete (idempotent)

Example create:
`POST /v1/publishers/penguin/books`
`{ "title": "Snow Crash", "isbn": "9780553380958", "language": "en" }`
→ `201 { "name": "publishers/penguin/books/bk_8f2e", "title": ... }`

## 3. The hard part — checkout & return
Custom methods on the copy (AIP-136), because lending is a state
transition with business rules, not resource editing:
- `POST /v1/copies/{copy}:checkout  { "member": "members/m_123" }`
  → `200` the loan record; `409 ALREADY_EXISTS` if out;
  `422/FAILED_PRECONDITION` if the member is at their limit.
- `POST /v1/copies/{copy}:return` → `200`; idempotent by design
  (returning an already-returned copy is a no-op success).
Each checkout also creates `loans/{loan}` (member, copy, due_date) —
the queryable history; `GET /v1/members/{member}/loans` lists it.
The pure-CRUD alternative (Create on a loans collection) is defensible;
what fails the bar is `PATCH copy.status` — it hides the rules.

## 4. Evolution & operations
Additive changes only: new optional fields, new methods. Never rename
or re-type a field (AIP-180). Errors share one envelope
`{ "error": { "code", "message", "details" } }` mapped to a small
fixed code set. Catalog reads get generous rate limits; lending
writes are per-branch quota'd. Version in the path (`/v1/`).
