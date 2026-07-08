class LRUCache:
    """Fixed-capacity cache that evicts the least recently used entry
    when full."""

    def __init__(self, capacity: int):
        self.capacity = capacity

    def get(self, key: int) -> int:
        """Return the value for key, or -1 if not present. Accessing a
        key marks it as most recently used."""
        raise NotImplementedError

    def put(self, key: int, value: int) -> None:
        """Insert or update key with value, marking it most recently
        used. If inserting a new key would exceed capacity, evict the
        least recently used entry first."""
        raise NotImplementedError
