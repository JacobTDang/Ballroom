import time


class TokenBucket:
    """A rate limiter shared by many threads. allow() takes a token if
    available; refill(n) (called by an external ticker) adds tokens,
    clamped at capacity.

    TODO: check-then-decrement below is two separate steps -- under
    contention this hands out more tokens than exist.
    """

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.tokens = capacity

    def allow(self) -> bool:
        if self.tokens > 0:
            time.sleep(0.0005)  # bookkeeping -- widens the race window
            self.tokens -= 1
            return True
        return False

    def refill(self, n: int):
        self.tokens = min(self.tokens + n, self.capacity)
