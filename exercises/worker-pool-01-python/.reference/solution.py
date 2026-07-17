import threading


def process_all(jobs, workers, fn):
    """A shared index cursor feeds worker threads; each writes
    results[i] for the i it claimed -- distinct slots, so only the
    cursor needs the lock. Position, not completion time, decides where
    a result lands, so ordering is free."""
    results = [None] * len(jobs)
    next_index = 0
    lock = threading.Lock()

    def worker():
        nonlocal next_index
        while True:
            with lock:
                i = next_index
                if i >= len(jobs):
                    return
                next_index += 1
            results[i] = fn(jobs[i])

    threads = [threading.Thread(target=worker) for _ in range(workers)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    return results
