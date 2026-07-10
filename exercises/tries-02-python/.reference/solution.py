class WordDictionary:
    """Supports adding words and searching, where a search query may
    use '.' to match any single character."""

    def __init__(self):
        self.children: dict[str, WordDictionary] = {}
        self.is_end = False

    def add_word(self, word: str) -> None:
        node = self
        for c in word:
            if c not in node.children:
                node.children[c] = WordDictionary()
            node = node.children[c]
        node.is_end = True

    def search(self, word: str) -> bool:
        return self._search_from(word, 0)

    def _search_from(self, word: str, idx: int) -> bool:
        node = self
        for i in range(idx, len(word)):
            c = word[i]
            if c == ".":
                return any(
                    child._search_from(word, i + 1) for child in node.children.values()
                )
            if c not in node.children:
                return False
            node = node.children[c]
        return node.is_end
