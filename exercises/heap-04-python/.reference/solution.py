import heapq


def find_kth_largest(nums: list[int], k: int) -> int:
    """Return the kth largest element of nums (1st largest is the
    maximum), via a min-heap capped at size k."""
    heap: list[int] = []
    for n in nums:
        heapq.heappush(heap, n)
        if len(heap) > k:
            heapq.heappop(heap)
    return heap[0]
