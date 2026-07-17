# LRU Cache With TTL

An LRU cache where entries *also* expire: capacity eviction (least
recently used goes first) and time expiry (an entry dies `ttl`
milliseconds after it was written) have to work together.

Time is injected (every operation takes `nowMillis`), so the tests are
exact. The rules:

- `Get` refreshes **recency** (LRU order) but never extends **expiry**
  — TTL runs from the write.
- Expired entries are misses, and they don't occupy capacity: a full
  cache whose entries are expired accepts new writes without evicting
  anything live.
- Writing an existing key updates its value, recency, and write time.

## The invariant the tests enforce

Exactly the rules above, each with a dedicated case — including the
interaction traps (expired-but-recently-used entries are still dead;
eviction skips corpses before touching live entries).

API: `TTLCache(capacity, ttl_ms)`, `.put_at(key, value, now_ms)`, `.get_at(key, now_ms) -> value | None`. `functools.lru_cache` and `OrderedDict.move_to_end` are off-limits — the mechanics are the exercise.
