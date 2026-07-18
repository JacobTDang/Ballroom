# API Design Roadmap

Twelve API-design questions practiced in Ballroom's API Design
category — each as a guided **coach** session (the 4-step method, one
step at a time, 90 minutes) and a timed **interviewer** mock (bare
prompt, you scope it, 45 minutes), graded on submit against hidden
rubrics grounded in the golden-standard references: **Google's API
Design Guide** (cloud.google.com/apis/design) and the **AIPs**
(google.aip.dev). Read the design guide once before Q1; read the AIP a
question cites after you've attempted it, not before — the gap between
your answer and the AIP is the lesson.

**Cadence**: one coach session per week, its interviewer mock one to
two weeks later (the picker's *mock due* marker tracks this). The
Fundamentals questions compound — Q1's resource model is the substrate
every later question builds on.

## How to practice (start here)

1. **Launch**: `ballroom` → Enter past the boot checks → `1` Practice →
   **API Design** → pick the question → **coach** first, interviewer
   when the mock comes due.
2. **Read the problem pane** (`M-1`): coach prompts walk the scope;
   interviewer prompts are one line on purpose — scoping is graded.
3. **Write the spec in `solution.md`** (`M-1`, the editor): resources,
   methods, wire shapes, and the question's hard mechanism. Concrete
   beats complete: real field names, real status codes, one worked
   example request/response.
4. **Talk to the tutor** (`M-2`): the coach walks the 4 steps; the
   interviewer probes like the real thing and answers scoping
   questions tersely.
5. **Submit with `M-q`**: the rubric and a distilled reference design
   land in your workspace, and the grader scores each dimension —
   compare structure and decisions, not wording.

**The method** — for every question, in order:

1. **Scope & resource model** — nouns, relationships, resource names
   (AIP-121/122).
2. **Methods & wire shapes** — standard methods before custom;
   request/response bodies; HTTP mapping and status codes
   (AIP-131–136).
3. **The hard part** — the question's featured mechanism, designed to
   the AIP bar, with a worked example.
4. **Evolution & operations** — what changes safely, what breaks
   clients, limits and operational reality (AIP-180/185).

## Phase 1 — Fundamentals

- [ ] **Design the Library API** (`library-api-01`) — resource-oriented
      CRUD, hierarchical names, and where checkout/return refuse to be
      CRUD (custom methods). AIP-121/122/131–136.
- [ ] **Paginate the Product Catalog** (`pagination-01`) — opaque page
      tokens over 100M rows; why offset pagination dies at scale;
      stability while the data moves. AIP-158/132.
- [ ] **Filter and Order the Fleet** (`list-filtering-01`) — a bounded
      filter grammar, order_by, and index honesty; search as a
      different thing. AIP-160/132.
- [ ] **Error Surface for a Payments API** (`error-model-01`) —
      canonical codes, machine-readable details, retryability; what
      never leaks. AIP-193.

## Phase 2 — Advanced mechanics

- [ ] **Transcode Video: Long-Running Operations**
      (`long-running-ops-01`) — the operation resource, progress,
      cancellation, result delivery. AIP-151.
- [ ] **Upload 10 GB, Resumably** (`resumable-upload-01`) — session
      initiation, offset negotiation, integrity, the small-file fast
      path. GCS resumable protocol, AIP-133.
- [ ] **Bulk Import for Contacts** (`batch-operations-01`) — batch
      methods, atomicity choices, per-item errors, limits.
      AIP-231–235.
- [ ] **Two Admins, One Settings Page** (`etag-concurrency-01`) —
      ETags, If-Match, 412; optimistic concurrency without locks.
      AIP-154.
- [ ] **Charge Exactly Once** (`retry-idempotency-01`) — idempotency
      keys, replay semantics, in-flight duplicates, fingerprint
      conflicts. AIP-155.
- [ ] **Evolve Without Breaking Anyone** (`api-versioning-01`) — the
      breaking-change taxonomy, wire compatibility, deprecation
      lifecycle. AIP-180/181/185.
- [ ] **Order Events: the Webhook Surface** (`webhooks-01`) —
      subscriptions as resources, at-least-once delivery, signing,
      retries, dead letters.
- [ ] **Rate Limits & Quotas for the Platform** (`api-quotas-01`) —
      quota vs rate limit, limit dimensions, the 429 contract,
      fairness.

## Build the mechanics

The implementation ladder's API-mechanics section is this roadmap's
hands-on half — design it here, build it there:
`cursor-pagination-01` (Q2's tokens), `idempotency-store-01` (Q9's
store), `conditional-request-01` (Q8's state machine),
`field-mask-update-01` (the AIP-134 update engine). See
`docs/implementation-roadmap.md`.

## Sibling roadmaps

`docs/system-design-roadmap.md` designs the systems these APIs front —
the two tracks share the coach/interviewer rhythm and the rubric bar.
`docs/implementation-roadmap.md` builds the mechanics; 
`docs/concurrency-roadmap.md` is what happens inside the handlers.
