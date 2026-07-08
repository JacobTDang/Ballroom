import threading
import time


class Counter:
    """Increments a shared counter across many threads. The final count
    should equal the number of increment() calls."""

    def __init__(self):
        self.count = 0

    def increment(self):
        current = self.count
        time.sleep(0)  # yield control here — exposes the race window
        self.count = current + 1


def run(n: int) -> int:
    c = Counter()
    threads = [threading.Thread(target=c.increment) for _ in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    return c.count
