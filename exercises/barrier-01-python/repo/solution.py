import threading


class Barrier:
    """Makes n threads rendezvous: each wait() blocks until all n have
    arrived, then all proceed -- and the barrier must be reusable for
    the next round.

    TODO: releasing via a one-shot Event works exactly once -- the
    second round finds it already set and nobody waits at all.
    """

    def __init__(self, n: int):
        self.n = n
        self.arrived = 0
        self.lock = threading.Lock()
        self.release = threading.Event()

    def wait(self):
        with self.lock:
            self.arrived += 1
            if self.arrived == self.n:
                self.arrived = 0
                self.release.set()
                return
        self.release.wait()
