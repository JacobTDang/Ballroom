from collections import deque


class SlidingWindow:
    """Keeps the timestamps of allowed requests (a deque: they're in
    order, so eviction pops from the left) and evicts entries a full
    window old before deciding. Denied requests are never recorded."""

    def __init__(self, limit: int, window_ms: int):
        self.limit = limit
        self.window_ms = window_ms
        self.allowed = deque()

    def allow_at(self, now_ms: int) -> bool:
        cutoff = now_ms - self.window_ms
        while self.allowed and self.allowed[0] <= cutoff:
            self.allowed.popleft()
        if len(self.allowed) < self.limit:
            self.allowed.append(now_ms)
            return True
        return False
