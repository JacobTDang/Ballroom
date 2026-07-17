import threading


class Lazy:
    """Double-checked locking, done right: the fast path reads the
    flag without the lock, the slow path re-checks under it -- so after
    the first call the lock is never touched again, and during it only
    one thread runs init."""

    def __init__(self, init):
        self._init = init
        self._value = None
        self._done = False
        self._lock = threading.Lock()

    def get(self):
        if not self._done:
            with self._lock:
                if not self._done:
                    self._value = self._init()
                    self._done = True
        return self._value
