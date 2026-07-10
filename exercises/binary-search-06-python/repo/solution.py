class TimeMap:
    """Stores multiple values per key, each tagged with the timestamp
    it was set at."""

    def __init__(self):
        pass

    def set(self, key: str, value: str, timestamp: int) -> None:
        raise NotImplementedError

    def get(self, key: str, timestamp: int) -> str:
        raise NotImplementedError
