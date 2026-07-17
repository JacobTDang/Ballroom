import threading


class TokenBucket:
    """One lock makes check-and-take a single atomic step -- the whole
    fix. refill clamps under the same lock so it can't race allow into
    overshooting capacity."""

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.tokens = capacity
        self.lock = threading.Lock()

    def allow(self) -> bool:
        with self.lock:
            if self.tokens > 0:
                self.tokens -= 1
                return True
            return False

    def refill(self, n: int):
        with self.lock:
            self.tokens = min(self.tokens + n, self.capacity)
