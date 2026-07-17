import threading


def fan_out_in(inputs, workers, stage):
    """Fan inputs out to `workers` threads running stage; collect every
    result (order doesn't matter).

    TODO: this version returns without joining the workers -- results
    go missing whenever a stage call is still running.
    """
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
            results.append(stage(inputs[i]))

    threads = [threading.Thread(target=worker) for _ in range(workers)]
    for t in threads:
        t.start()
    return results
