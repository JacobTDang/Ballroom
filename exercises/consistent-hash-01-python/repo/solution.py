class Ring:
    """Maps keys to nodes. Adding/removing a node should only remap
    the keys in its neighborhood.

    TODO: hash(key) % len(nodes) remaps almost EVERY key whenever the
    node count changes -- the exact failure consistent hashing exists
    to fix. (vnodes is ignored here too.)
    """

    def __init__(self, vnodes: int):
        self.nodes = []

    def _hash(self, key: str) -> int:
        h = 0xCBF29CE484222325
        for b in key.encode():
            h ^= b
            h = (h * 0x100000001B3) % (1 << 64)
        return h

    def add_node(self, name: str):
        self.nodes.append(name)
        self.nodes.sort()

    def remove_node(self, name: str):
        if name in self.nodes:
            self.nodes.remove(name)

    def lookup(self, key: str) -> str:
        if not self.nodes:
            return ""
        return self.nodes[self._hash(key) % len(self.nodes)]
