import threading


def run_limited(tasks, limit):
    """Run every task, with at most `limit` executing concurrently.

    TODO: this version launches everything at once -- the limit is
    ignored entirely.
    """
    threads = [threading.Thread(target=t) for t in tasks]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
