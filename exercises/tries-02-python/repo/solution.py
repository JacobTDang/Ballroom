class WordDictionary:
    """Supports adding words and searching, where a search query may
    use '.' to match any single character."""

    def __init__(self):
        pass

    def add_word(self, word: str) -> None:
        raise NotImplementedError

    def search(self, word: str) -> bool:
        raise NotImplementedError
