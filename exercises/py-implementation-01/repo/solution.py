import time


class RateLimiter:
    """Fixed-window rate limiter: allows at most `limit` calls per
    `window_seconds`. All calls within the same window count toward the
    same limit; the window resets `window_seconds` after the first call
    in it."""

    def __init__(self, limit: int, window_seconds: float):
        self.limit = limit
        self.window_seconds = window_seconds
        # TODO: implement

    def allow(self) -> bool:
        """Return True if a new request should be allowed right now."""
        raise NotImplementedError
