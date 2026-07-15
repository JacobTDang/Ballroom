# Reference Design — Pastebin / Bit.ly

A solid answer, distilled. Compare structure and decisions, not wording.

## 1. Use cases, constraints, estimates
Shorten + redirect in scope; analytics minimal (click count), custom
aliases out. 100M links/month → **~40 writes/s**; 10:1 reads →
**~400 reads/s**. Record ~500 B → 500 B × 100M/mo × 36 mo ≈ **1.8 TB**.
Read-heavy: design the read path first.

## 2. High-level design
Write: client → LB → write API → ID generation → SQL store.
Read: client → LB → read API → cache → store → **301/302 redirect**.

## 3. Core components
- **Code generation:** base-62 of an auto-increment ID — collision-free
  by construction; 62^7 ≈ 3.5T covers decades. (Alternative: MD5(url +
  salt) truncated + collision retry — argue either, handle collisions.)
- **Storage:** a single `links` table (short_code PK, url, created_at,
  expires_at, click_count). SQL is fine: one key-lookup access pattern,
  1.8 TB fits standard tooling. NoSQL KV equally defensible.
- **Redirect semantics:** 301 = browser caches → fewer hits but blind
  analytics; 302 = every click observed. Pick one and say why.
- **API:** `POST /links {url}` → `{code}`; `GET /{code}` → redirect.

## 4. Scale
- Cache popular codes (cache-aside, LRU): 20% of a month's links ≈
  10 GB — one Redis box covers most reads.
- Read replicas next; shard by hash(code) only when writes outgrow one
  primary.
- Expiry via lazy delete on read + periodic sweep.
- Click counts buffered/batched, never synchronous row UPDATEs.
