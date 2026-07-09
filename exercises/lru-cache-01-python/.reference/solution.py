from collections import OrderedDict


class LRUCache:
    """Fixed-capacity cache that evicts the least recently used entry
    when full."""

    def __init__(self, capacity: int):
        self.capacity = capacity
        self._data: "OrderedDict[int, int]" = OrderedDict()

    def get(self, key: int) -> int:
        """Return the value for key, or -1 if not present. Accessing a
        key marks it as most recently used."""
        if key not in self._data:
            return -1
        self._data.move_to_end(key)
        return self._data[key]

    def put(self, key: int, value: int) -> None:
        """Insert or update key with value, marking it most recently
        used. If inserting a new key would exceed capacity, evict the
        least recently used entry first."""
        if key in self._data:
            self._data.move_to_end(key)
        self._data[key] = value
        if len(self._data) > self.capacity:
            self._data.popitem(last=False)
