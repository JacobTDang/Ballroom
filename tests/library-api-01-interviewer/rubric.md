# Grading Rubric — Design the Library API

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Resource model & naming

- Resources are nouns, not verbs — no `/getBook` or `/doCheckout`
  endpoints anywhere (AIP-121).
- Concrete hierarchical names in AIP-122 shape
  (`publishers/{publisher}/books/{book}`), with the hierarchy justified
  from ownership/use — or a flat model with the trade-off argued.

## 2. Standard-method coverage

- The five standard methods (AIP-131..135) mapped to HTTP correctly:
  GET resource, GET collection, POST on the collection, PATCH, DELETE.
- Idempotency of each stated (GET/PUT/DELETE yes, POST no, PATCH
  depends) — not just assumed.

## 3. Custom-method judgment

- Checkout/return modeled deliberately: `POST .../copies/{copy}:checkout`
  custom methods (AIP-136) or a loans collection — either accepted,
  but the trade-off must be argued, not bent into fake CRUD like
  `PATCH copy.status = "out"`.
- Failure cases specified: already checked out, member at limit.

## 4. Request/response shapes

- Bodies are the resource itself; Create returns the created resource
  with its server-assigned name.
- Real field names with deliberate types/formats (timestamps, enums),
  at least one full worked request/response example.

## 5. List basics

- Collection List sketched with `page_size`/`page_token` (AIP-158) —
  full pagination depth is the next question's job, but its absence
  here is a miss.

## 6. Error surface

- 404 vs 400 vs 409 assigned to the obvious cases (missing book,
  malformed request, checkout conflict) consistently.
- One error envelope shape shared by every endpoint.

## 7. Communication & trade-offs

- REST-vs-RPC and hierarchy-vs-flat argued from actual use cases.
- The candidate drove the 4 steps themselves; decisions carry reasons,
  not just conclusions.
