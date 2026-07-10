class Trie:
    """Prefix tree over lowercase English letters."""

    def __init__(self):
        pass

    def insert(self, word: str) -> None:
        raise NotImplementedError

    def search(self, word: str) -> bool:
        raise NotImplementedError

    def starts_with(self, prefix: str) -> bool:
        raise NotImplementedError
