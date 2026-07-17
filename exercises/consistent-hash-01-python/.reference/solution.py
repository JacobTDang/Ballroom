import bisect


def _fnv1a(data: bytes) -> int:
    h = 0xCBF29CE484222325
    for b in data:
        h ^= b
        h = (h * 0x100000001B3) % (1 << 64)
    return h


class Ring:
    """Each node occupies `vnodes` positions on the ring; a key
    belongs to the first position clockwise from its hash (bisect over
    the sorted positions). Removing a node deletes only its own
    positions, so every other key keeps its owner."""

    def __init__(self, vnodes: int):
        self.vnodes = vnodes
        self.positions = []
        self.owner = {}

    def add_node(self, name: str):
        for i in range(self.vnodes):
            p = _fnv1a(f"{name}#{i}".encode())
            if p in self.owner:
                continue  # vanishing collision odds: first owner keeps it
            self.owner[p] = name
            bisect.insort(self.positions, p)

    def remove_node(self, name: str):
        self.positions = [p for p in self.positions if self.owner.get(p) != name]
        self.owner = {p: n for p, n in self.owner.items() if n != name}

    def lookup(self, key: str) -> str:
        if not self.positions:
            return ""
        h = _fnv1a(key.encode())
        i = bisect.bisect_left(self.positions, h)
        if i == len(self.positions):
            i = 0  # wrap around the ring
        return self.owner[self.positions[i]]
