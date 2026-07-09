import time


class RateLimiter:
    """Fixed-window rate limiter: allows at most `limit` calls per
    `window_seconds`. All calls within the same window count toward the
    same limit; the window resets `window_seconds` after the first call
    in it."""

    def __init__(self, limit: int, window_seconds: float):
        self.limit = limit
        self.window_seconds = window_seconds
        self._window_start = None
        self._count = 0

    def allow(self) -> bool:
        """Return True if a new request should be allowed right now."""
        now = time.time()
        if self._window_start is None or now - self._window_start >= self.window_seconds:
            self._window_start = now
            self._count = 0
        if self._count >= self.limit:
            return False
        self._count += 1
        return True
