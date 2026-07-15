# Design a Web Crawler

Design a distributed web crawler: given seed URLs, it fetches pages,
extracts links, and keeps crawling — feeding page content to a search
index while staying polite to the sites it visits.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: crawl from seeds, extract + follow links, store page
  content for indexing, re-crawl for freshness. Out of scope? (Ranking,
  serving search queries, JavaScript rendering.)
- The two correctness problems every crawler must solve: not fetching
  the same URL twice, and not hammering one site.
- Put numbers on it: pages/month, fetches/sec, average page size,
  storage growth, how often pages are refreshed.

## Suggested defaults (if you want a starting point)

- 1 billion pages crawled per month, refreshed on a weekly cadence
- Average page ~500 KB HTML; store extracted text, not raw bytes
- Politeness: at most 1 request/sec per domain, honor robots.txt
- Crawl order should prefer important/fresh pages, but simple FIFO is
  an acceptable starting point

## What good looks like

By the end you should have: fetch-rate arithmetic from the monthly
target, a high-level loop (frontier → fetcher → parser → dedup →
frontier), a concrete dedup answer for both URLs and content, a
politeness mechanism that actually serializes per-domain requests,
crawler-trap defenses, and a scaling story for distributing the
frontier across machines.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
