import queue
import threading


class Server:
    """Runs handle on submitted jobs with a pool of worker threads.
    stop() must drain everything accepted, refuse new work, and only
    return when the workers are done.

    TODO: this stop flips the flag and returns immediately -- queued
    jobs are abandoned and the workers are left running (daemon
    threads die with the process, taking unfinished work with them).
    """

    def __init__(self, workers, handle):
        self.jobs = queue.Queue()
        self.stopped = False
        self.threads = []
        for _ in range(workers):
            t = threading.Thread(target=self._worker, args=(handle,), daemon=True)
            t.start()
            self.threads.append(t)

    def _worker(self, handle):
        while True:
            v = self.jobs.get()
            handle(v)

    def submit(self, v) -> bool:
        if self.stopped:
            return False
        self.jobs.put(v)
        return True

    def stop(self):
        self.stopped = True
