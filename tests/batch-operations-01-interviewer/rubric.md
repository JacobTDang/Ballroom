# Grading Rubric — Bulk Import for Contacts

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Batch method shape

- Batch methods mirror their singular methods — requests embed the
  singular request shapes with the same names and fields (AIP-233
  style), not a bespoke import format.
- Custom-method spelling on the collection (e.g. `:batchCreate`) with
  sane HTTP mapping.

## 2. Atomicity decision

- An explicit choice per method: batchGet atomic (AIP-231); for
  mutations, all-or-nothing vs partial success chosen and argued from
  the import use case — not left implicit.
- The choice is visible in the response shape (a failed atomic batch
  returns one error; partial success returns per-item results).

## 3. Per-item error reporting

- Results are position-matched to the request array (or keyed
  equivalently); every input slot gets a resource or a status —
  nothing silently dropped.
- Item errors reuse the canonical error model (AIP-193), not
  free-text strings.

## 4. Limits

- Concrete limits on items per batch and payload bytes, with a
  rationale tied to the numbers; over-limit returns INVALID_ARGUMENT
  naming the limit.
- Client guidance for chunking larger imports.

## 5. Idempotency of retried batches

- A half-succeeded batch retried does not duplicate contacts: per-item
  request ids / client-assigned ids (AIP-155/133) or an equivalent
  dedup story, stated concretely.
- The retry story covers both partial-success and
  failed-atomic-batch cases.

## 6. Throughput reality

- Batch vs singular quantified (round trips, payload overhead) at the
  stated import sizes.
- The ceiling is named: at some size a synchronous batch stops making
  sense and an asynchronous import job (operation) takes over.

## 7. Communication & trade-offs

- Atomicity-vs-availability and simplicity-vs-throughput trade-offs
  argued; limits and choices tied to the step-1 numbers; the
  candidate drove the 4 steps themselves.
