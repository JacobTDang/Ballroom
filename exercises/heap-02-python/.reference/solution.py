import heapq


def last_stone_weight(stones: list[int]) -> int:
    """Repeatedly smash the two heaviest stones together and return
    the weight of whatever stone (if any) remains."""
    heap = [-s for s in stones]
    heapq.heapify(heap)
    while len(heap) > 1:
        a = -heapq.heappop(heap)
        b = -heapq.heappop(heap)
        if a != b:
            heapq.heappush(heap, -(a - b))
    return -heap[0] if heap else 0
