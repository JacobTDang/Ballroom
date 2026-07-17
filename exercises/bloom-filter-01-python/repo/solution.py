class BloomFilter:
    """A bit array + k hashes. add() sets k bits; might_contain()
    checks them -- "definitely not" or "probably yes".

    TODO: one weak hash into a fixed 64-slot table, ignoring both
    parameters -- no false negatives, but the table saturates instantly
    and almost every absent key collides.
    """

    def __init__(self, bits: int, hashes: int):
        self.table = [False] * 64

    def _hash(self, key: str) -> int:
        return sum(ord(c) for c in key) % 64

    def add(self, key: str):
        self.table[self._hash(key)] = True

    def might_contain(self, key: str) -> bool:
        return self.table[self._hash(key)]
