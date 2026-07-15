# Grading Rubric — Web Crawler

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- Scoped to crawling + content storage for indexing, with explicit
  out-of-scope calls (ranking, query serving, JS rendering).
- Named both correctness problems unprompted: duplicate work and
  politeness.

## 2. Back-of-envelope estimates

- Fetches/sec derived from the monthly page target, arithmetic shown.
- Storage sized from page count × stored-content size, over the
  refresh window.

## 3. High-level design

- The crawl loop drawn as components: frontier → fetchers → parser →
  dedup → storage, with extracted links feeding back into the frontier.
- Search indexing shown as a downstream consumer, not conflated with
  the crawler itself.

## 4. Deduplication

- URL-seen test that works at scale (hash set with sharding, or bloom
  filter with its false-positive trade-off stated).
- Content-level dedup addressed (page signature / simhash) with what
  happens on a match.

## 5. Politeness

- A mechanism that actually serializes per-domain fetches (per-domain
  queues owned by one worker, or a domain-keyed lock/token bucket) —
  not just "we rate limit".
- robots.txt fetched, cached, and honored.

## 6. Crawler traps & robustness

- Defenses named: max depth, max URLs per domain, URL normalization,
  detecting infinite/generated link spaces.
- Fetch failures and retries handled without poisoning the frontier.

## 7. Scaling story

- Frontier partitioned by domain (keeps politeness local to one node),
  distributed via consistent hashing; rebalancing on node loss.
- Prioritization/freshness discussed (importance queues or re-crawl
  scheduling).

## 8. Communication & trade-offs

- Trade-offs stated (bloom filter vs exact set, breadth vs depth,
  freshness vs throughput), driven by the step-1 numbers; the candidate
  drove the 4-step structure.
