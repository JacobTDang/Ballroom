import heapq


class MedianFinder:
    """Tracks the running median of a stream of integers, using two
    heaps that split the stream in half: small holds the lower half
    (as negated values, so Python's min-heap acts as a max-heap) and
    large holds the upper half (min-heap). Kept balanced within 1 of
    each other after every insert, so the median is always at the
    top of one (or both) heaps."""

    def __init__(self):
        self.small: list[int] = []  # max-heap, stored negated
        self.large: list[int] = []  # min-heap

    def add_num(self, num: int) -> None:
        heapq.heappush(self.small, -num)
        heapq.heappush(self.large, -heapq.heappop(self.small))
        if len(self.large) > len(self.small):
            heapq.heappush(self.small, -heapq.heappop(self.large))

    def find_median(self) -> float:
        if len(self.small) > len(self.large):
            return float(-self.small[0])
        return (-self.small[0] + self.large[0]) / 2.0
