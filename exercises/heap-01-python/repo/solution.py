class KthLargest:
    """Tracks the kth largest value seen so far in a stream of
    integers."""

    def __init__(self, k: int, nums: list[int]):
        self.k = k

    def add(self, val: int) -> int:
        raise NotImplementedError
