import time


class Lazy:
    """Computes a value on first use -- the init function is expensive
    and must run exactly once, no matter how many threads call get()
    concurrently.

    TODO: the check below isn't atomic with the assignment -- two
    threads can both see done == False and both run init.
    """

    def __init__(self, init):
        self._init = init
        self._value = None
        self._done = False

    def get(self):
        if not self._done:
            time.sleep(0)  # yield -- widens the double-init window
            self._value = self._init()
            self._done = True
        return self._value
