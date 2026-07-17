# Bloom Filter

A space-efficient set that can say "definitely not present" or
"probably present": a bit array plus k hash functions. Add sets k bits;
membership checks them.

The starter uses one weak hash into a tiny table — its false-positive
rate is catastrophic. Yours gets a real bit array and k derived hashes
(double hashing `h1 + i*h2` over one or two base hashes is the
standard trick — no need for k independent hash functions).

## The invariant the tests enforce

- **Zero false negatives** — everything added is always reported
  present. This is the bloom filter's hard guarantee.
- **Bounded false positives** — with 16384 bits, 4 hashes, and 500
  added keys, fewer than 2% of 10,000 absent keys may report present
  (theory says ~0.02% — the budget is generous).
- An empty filter reports nothing present.

API: `class BloomFilter { BloomFilter(int bits, int hashes); void Add(const std::string&); bool MightContain(const std::string&); }`.

Think: why does the false-positive rate collapse when bits are too few
— and why can a bloom filter never have a false negative?
