class MyHashMap:
    """Hash map for non-negative integer keys and values, built without
    the language's own dict."""

    def __init__(self):
        pass

    def put(self, key: int, value: int) -> None:
        """Insert key with value, or update it if key already exists."""
        raise NotImplementedError

    def get(self, key: int) -> int:
        """Return the value for key, or -1 if key is absent."""
        raise NotImplementedError

    def remove(self, key: int) -> None:
        """Delete key if present; do nothing otherwise."""
        raise NotImplementedError
