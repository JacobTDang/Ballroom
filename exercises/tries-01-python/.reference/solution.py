class Trie:
    """Prefix tree over lowercase English letters."""

    def __init__(self):
        self.children: dict[str, Trie] = {}
        self.is_end = False

    def insert(self, word: str) -> None:
        node = self
        for c in word:
            if c not in node.children:
                node.children[c] = Trie()
            node = node.children[c]
        node.is_end = True

    def search(self, word: str) -> bool:
        node = self._find(word)
        return node is not None and node.is_end

    def starts_with(self, prefix: str) -> bool:
        return self._find(prefix) is not None

    def _find(self, s: str) -> "Trie | None":
        node = self
        for c in s:
            if c not in node.children:
                return None
            node = node.children[c]
        return node
