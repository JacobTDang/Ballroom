import heapq


class KthLargest:
    """Tracks the kth largest value seen so far in a stream of
    integers, using a min-heap capped at size k -- the heap's
    smallest element (index 0) is always the kth largest overall."""

    def __init__(self, k: int, nums: list[int]):
        self.k = k
        self.heap = list(nums)
        heapq.heapify(self.heap)
        while len(self.heap) > k:
            heapq.heappop(self.heap)

    def add(self, val: int) -> int:
        heapq.heappush(self.heap, val)
        if len(self.heap) > self.k:
            heapq.heappop(self.heap)
        return self.heap[0]
