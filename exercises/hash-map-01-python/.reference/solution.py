class MyHashMap:
    """Hash map for non-negative integer keys and values, built without
    the language's own dict."""

    _BUCKETS = 1024

    def __init__(self):
        self._buckets = [[] for _ in range(self._BUCKETS)]

    def _bucket(self, key: int) -> list:
        return self._buckets[key % self._BUCKETS]

    def put(self, key: int, value: int) -> None:
        bucket = self._bucket(key)
        for i, (k, _) in enumerate(bucket):
            if k == key:
                bucket[i] = (key, value)
                return
        bucket.append((key, value))

    def get(self, key: int) -> int:
        for k, v in self._bucket(key):
            if k == key:
                return v
        return -1

    def remove(self, key: int) -> None:
        bucket = self._bucket(key)
        for i, (k, _) in enumerate(bucket):
            if k == key:
                bucket.pop(i)
                return
