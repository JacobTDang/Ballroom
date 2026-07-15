# Reference Design — Web Crawler

## 1. Use cases, constraints, estimates
1B pages/month → **~400 fetches/s** sustained. ~500 KB/page fetched,
store extracted text (~10 KB) → ~10 TB/month of text. Weekly refresh
doubles effective load. Politeness: ≤1 req/s per domain.

## 2. High-level design
The loop: **URL frontier → fetchers → parser/extractor → dedup →
content store**, with extracted links normalized and fed back into the
frontier. Search indexing consumes the content store downstream.

## 3. Core components
- **URL dedup:** "seen" set at billions of URLs — sharded hash set, or
  a bloom filter with its false-positive trade-off stated (a missed
  crawl is acceptable; a wasted fetch is not fatal).
- **Content dedup:** page signature (hash or simhash for near-dupes);
  on match, record the alias, don't re-store or re-extract links.
- **Politeness:** partition the frontier **by domain** so exactly one
  worker owns a domain's queue — serialization is structural, not
  lock-based. Cache robots.txt per domain with TTL.
- **Traps:** max depth, max URLs/domain, URL normalization (strip
  session ids), cycle detection via the seen set.

## 4. Scale
- Frontier partitioned by hash(domain) across nodes (consistent
  hashing); a dead node's domains reassign; in-flight URLs recovered
  from a checkpoint log.
- Priority queues layered on: importance and refresh-due ordering.
- Fetchers scale horizontally; DNS caching matters at 400 req/s.
