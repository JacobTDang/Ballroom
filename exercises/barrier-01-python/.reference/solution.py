import threading


class Barrier:
    """A generation counter makes reuse safe: each round waits on its
    own generation, so a round-2 arrival can never consume a round-1
    release. The last arriver flips the generation and notifies all;
    waiters loop until THEIR generation has passed."""

    def __init__(self, n: int):
        self.n = n
        self.arrived = 0
        self.generation = 0
        self.cond = threading.Condition()

    def wait(self):
        with self.cond:
            gen = self.generation
            self.arrived += 1
            if self.arrived == self.n:
                self.arrived = 0
                self.generation += 1
                self.cond.notify_all()
                return
            while gen == self.generation:
                self.cond.wait()
