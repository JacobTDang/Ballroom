import heapq


def min_interval(intervals: list[list[int]], queries: list[int]) -> list[int]:
    """For each query, return the size of the smallest interval that
    contains it (left <= query <= right), or -1 if no interval contains
    it. The result is in the same order as queries."""
    sorted_intervals = sorted(intervals, key=lambda iv: iv[0])
    order = sorted(range(len(queries)), key=lambda i: queries[i])

    result = [0] * len(queries)
    heap: list[tuple[int, int]] = []  # (size, end)
    i = 0

    for idx in order:
        q = queries[idx]
        while i < len(sorted_intervals) and sorted_intervals[i][0] <= q:
            left, right = sorted_intervals[i]
            heapq.heappush(heap, (right - left + 1, right))
            i += 1

        while heap and heap[0][1] < q:
            heapq.heappop(heap)

        result[idx] = heap[0][0] if heap else -1

    return result
