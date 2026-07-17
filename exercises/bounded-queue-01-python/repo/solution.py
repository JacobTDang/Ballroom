import threading


class BoundedQueue:
    """A fixed-capacity FIFO shared between producer and consumer
    threads. put() must block while full; get() must block while empty.

    TODO: this version has neither bound nor blocking -- put ignores
    capacity and get returns None when empty.
    """

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.items = []
        self.lock = threading.Lock()

    def put(self, v):
        with self.lock:
            self.items.append(v)

    def get(self):
        with self.lock:
            if not self.items:
                return None
            return self.items.pop(0)
