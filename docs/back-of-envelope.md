# Back-of-Envelope Cheat Sheet

Inside a practice session: `less ~/back-of-envelope.md` in the terminal
pane (M-3).

## Powers of two

| Power | Exact       | Approx | Bytes    |
|-------|-------------|--------|----------|
| 10    | 1,024       | ~1 thousand | 1 KB |
| 20    | 1,048,576   | ~1 million  | 1 MB |
| 30    | ~1.07e9     | ~1 billion  | 1 GB |
| 40    | ~1.10e12    | ~1 trillion | 1 TB |
| 50    | ~1.13e15    | ~1 quadrillion | 1 PB |

## Latency numbers every programmer should know

| Operation                              | Time      | Scale it |
|----------------------------------------|-----------|----------|
| L1 cache reference                     | 0.5 ns    | |
| Branch mispredict                      | 5 ns      | |
| L2 cache reference                     | 7 ns      | 14× L1 |
| Mutex lock/unlock                      | 25 ns     | |
| Main memory reference                  | 100 ns    | 200× L1 |
| Compress 1 KB (Zippy)                  | 10 µs     | |
| Send 1 KB over 1 Gbps network          | 10 µs     | |
| Read 1 MB sequentially from memory     | 250 µs    | |
| Round trip within same datacenter      | 500 µs    | |
| Read 1 MB sequentially from SSD        | 1 ms      | 4× memory |
| Disk seek                              | 10 ms     | 20× datacenter RT |
| Read 1 MB sequentially from disk       | 30 ms     | 120× memory |
| Packet CA → Netherlands → CA           | 150 ms    | |

Rules of thumb that fall out of the table:

- Memory is fast, disk is slow, seeks are what kill you — read
  sequentially where possible.
- Compression is cheap relative to the network: compress before you
  send.
- Cross-datacenter round trips are ~300× in-datacenter ones — replicate
  data close to users.

## Time shortcuts

| Window            | Seconds (approx) |
|-------------------|------------------|
| 1 day             | ~86,400 (use 100K for margin) |
| 1 month           | ~2.5 million |
| 1 year            | ~31.5 million (use "π × 10^7") |

**The workhorse conversion: 1 million per month ≈ 0.4 per second.**
So 100M writes/month ≈ 40 writes/sec; 10B searches/month ≈ 4,000/sec.

## Standard estimation patterns

**QPS from monthly volume**
`X per month ÷ 2.5M sec ≈ X/2.5 per second (in millions)`
Peak ≈ 2–5× average — say which you're using.

**Storage**
`records × bytes/record × retention` — then sanity-check the record
size by listing its fields (ids ~8B, timestamps 8B, short strings
~tens of bytes, URLs ~hundreds).

**Bandwidth**
`QPS × payload size`. Egress usually dominates for read-heavy media
(video, images).

**Memory for a cache**
`hot fraction × total dataset` — the 80/20 rule: ~20% of objects serve
~80% of traffic. Then: does it fit on one box (~64–256 GB), or do you
need a cluster?

**Servers**
`peak QPS ÷ per-server capacity` — a commodity app server handles
~1K–10K simple req/s; say your assumption out loud.

## Worked micro-example (URL shortener)

- 100M new links/month → 100/2.5 ≈ **40 writes/sec**
- 10:1 reads → **400 reads/sec**
- 500 B/record × 100M/month × 36 months ≈ **1.8 TB** total
- Cache 20% of a month's links: 0.2 × 100M × 500 B ≈ **10 GB** — one box

## In an interview

1. Write the assumptions down before the arithmetic.
2. Round aggressively (2.5M sec/month, not 2,592,000).
3. Carry units through every line.
4. Sanity-check the result against something you know ("1.8 TB — fits
   on one disk, so storage isn't the hard part here").
