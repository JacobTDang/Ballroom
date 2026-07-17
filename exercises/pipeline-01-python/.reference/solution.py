import threading


def fan_out_in(inputs, workers, stage):
    """Same fan-out, but the caller joins every worker before
    returning -- and the append itself goes under the lock: list.append
    is atomic in CPython today, but the invariant shouldn't depend on
    that implementation detail."""
    results = []
    cursor = 0
    lock = threading.Lock()

    def worker():
        nonlocal cursor
        while True:
            with lock:
                i = cursor
                if i >= len(inputs):
                    return
                cursor += 1
            v = stage(inputs[i])
            with lock:
                results.append(v)

    threads = [threading.Thread(target=worker) for _ in range(workers)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    return results
