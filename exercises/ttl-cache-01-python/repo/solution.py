class TTLCache:
    """LRU eviction at capacity, plus per-entry expiry ttl_ms after
    the write. get_at returns the value or None.

    TODO: a plain dict -- no eviction, no expiry, no recency. Every
    rule in the problem statement is still yours to build.
    """

    def __init__(self, capacity: int, ttl_ms: int):
        self.capacity = capacity
        self.ttl_ms = ttl_ms
        self.items = {}

    def put_at(self, key, value, now_ms):
        self.items[key] = value

    def get_at(self, key, now_ms):
        return self.items.get(key)
