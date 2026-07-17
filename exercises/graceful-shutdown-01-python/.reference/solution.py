import queue
import threading

_SENTINEL = object()


class Server:
    """One sentinel per worker is the "no more work" signal: stop()
    marks stopped under the lock (so no submit sneaks in after), queues
    exactly `workers` sentinels behind the real jobs, and joins every
    worker -- FIFO order guarantees the sentinels are seen only after
    everything accepted has been handled."""

    def __init__(self, workers, handle):
        self.jobs = queue.Queue()
        self.stopped = False
        self.lock = threading.Lock()
        self.threads = []
        for _ in range(workers):
            t = threading.Thread(target=self._worker, args=(handle,))
            t.start()
            self.threads.append(t)

    def _worker(self, handle):
        while True:
            v = self.jobs.get()
            if v is _SENTINEL:
                return
            handle(v)

    def submit(self, v) -> bool:
        with self.lock:
            if self.stopped:
                return False
            self.jobs.put(v)
            return True

    def stop(self):
        with self.lock:
            if self.stopped:
                return
            self.stopped = True
            for _ in self.threads:
                self.jobs.put(_SENTINEL)
        for t in self.threads:
            t.join()
