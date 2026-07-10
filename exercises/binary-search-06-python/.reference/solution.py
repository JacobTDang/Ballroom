class TimeMap:
    """Stores multiple values per key, each tagged with the timestamp
    it was set at."""

    def __init__(self):
        self._store: dict[str, list[tuple[int, str]]] = {}

    def set(self, key: str, value: str, timestamp: int) -> None:
        # Timestamps arrive strictly increasing, so the per-key list
        # stays sorted without needing to insert.
        self._store.setdefault(key, []).append((timestamp, value))

    def get(self, key: str, timestamp: int) -> str:
        entries = self._store.get(key, [])
        lo, hi = 0, len(entries) - 1
        res = ""
        while lo <= hi:
            mid = lo + (hi - lo) // 2
            if entries[mid][0] <= timestamp:
                res = entries[mid][1]
                lo = mid + 1
            else:
                hi = mid - 1
        return res
