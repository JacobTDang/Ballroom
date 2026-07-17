class TTLCache:
    """Each entry remembers its write time (expiry never moves) and a
    monotonic touch sequence (recency does). Puts at capacity purge
    expired corpses first -- they don't deserve an eviction -- then
    evict the least recently used live entry. O(n) scans are fine at
    exercise scale; a production cache pairs the dict with a linked
    list."""

    def __init__(self, capacity: int, ttl_ms: int):
        self.capacity = capacity
        self.ttl_ms = ttl_ms
        self.items = {}  # key -> [value, written_at, touched_seq]
        self.seq = 0

    def _expired(self, entry, now_ms):
        return now_ms - entry[1] >= self.ttl_ms

    def put_at(self, key, value, now_ms):
        self.seq += 1
        if key in self.items:
            self.items[key] = [value, now_ms, self.seq]
            return
        for k in [k for k, e in self.items.items() if self._expired(e, now_ms)]:
            del self.items[k]
        if len(self.items) >= self.capacity:
            lru = min(self.items, key=lambda k: self.items[k][2])
            del self.items[lru]
        self.items[key] = [value, now_ms, self.seq]

    def get_at(self, key, now_ms):
        entry = self.items.get(key)
        if entry is None or self._expired(entry, now_ms):
            return None
        self.seq += 1
        entry[2] = self.seq  # recency refreshed; written_at deliberately not
        return entry[0]
