import threading


class Counter:
    """Increments a shared counter across many threads. The final count
    always equals the number of increment() calls."""

    def __init__(self):
        self.count = 0
        self._lock = threading.Lock()

    def increment(self):
        with self._lock:
            current = self.count
            self.count = current + 1


def run(n: int) -> int:
    c = Counter()
    threads = [threading.Thread(target=c.increment) for _ in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    return c.count
