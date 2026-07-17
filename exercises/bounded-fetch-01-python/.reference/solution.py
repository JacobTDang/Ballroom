import threading


def run_limited(tasks, limit):
    """The semaphore wraps the task body -- acquired inside the
    thread, immediately around the work -- so the bound covers
    execution, not thread creation."""
    sem = threading.Semaphore(limit)

    def bounded(task):
        with sem:
            task()

    threads = [threading.Thread(target=bounded, args=(t,)) for t in tasks]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
