import threading
from collections import deque


class BoundedQueue:
    """Condition-variable bounded FIFO: one lock, two wait conditions
    (not-full for producers, not-empty for consumers). wait() runs in a
    while-loop because wake-ups race with other threads taking the slot
    first."""

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.items = deque()
        self.lock = threading.Lock()
        self.not_full = threading.Condition(self.lock)
        self.not_empty = threading.Condition(self.lock)

    def put(self, v):
        with self.not_full:
            while len(self.items) >= self.capacity:
                self.not_full.wait()
            self.items.append(v)
            self.not_empty.notify()

    def get(self):
        with self.not_empty:
            while not self.items:
                self.not_empty.wait()
            v = self.items.popleft()
            self.not_full.notify()
            return v
