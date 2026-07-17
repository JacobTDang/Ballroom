class BloomFilter:
    """Double hashing: two FNV-1a variants generate all k probe
    positions as h1 + i*h2 (Kirsch-Mitzenmacher) -- k hashes as good as
    k independent ones without computing k real hashes."""

    def __init__(self, bits: int, hashes: int):
        self.bits = [False] * bits
        self.hashes = hashes

    @staticmethod
    def _fnv1a(data: bytes) -> int:
        h = 0xCBF29CE484222325
        for byte in data:
            h ^= byte
            h = (h * 0x100000001B3) % (1 << 64)
        return h

    @staticmethod
    def _mix(h: int) -> int:
        # splitmix64 finalizer: FNV alone clusters on similar keys over
        # power-of-two table sizes; the avalanche makes probes behave
        # independently.
        mask = (1 << 64) - 1
        h ^= h >> 30
        h = (h * 0xBF58476D1CE4E5B9) & mask
        h ^= h >> 27
        h = (h * 0x94D049BB133111EB) & mask
        h ^= h >> 31
        return h

    def _positions(self, key: str):
        h1 = self._mix(self._fnv1a(key.encode()))
        h2 = self._mix(h1 ^ 0x9E3779B97F4A7C15) | 1
        return [(h1 + i * h2) % len(self.bits) for i in range(self.hashes)]

    def add(self, key: str):
        for p in self._positions(key):
            self.bits[p] = True

    def might_contain(self, key: str) -> bool:
        return all(self.bits[p] for p in self._positions(key))
