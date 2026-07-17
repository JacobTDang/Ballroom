import threading
import time

from solution import Lazy


def test_init_runs_exactly_once_under_contention():
    calls = []
    lock = threading.Lock()

    def init():
        with lock:
            calls.append(1)
        time.sleep(0.01)
        return 42

    lazy = Lazy(init)
    results = [None] * 50

    def getter(i):
        results[i] = lazy.get()

    threads = [threading.Thread(target=getter, args=(i,)) for i in range(50)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(calls) == 1, f"init ran {len(calls)} times under contention, want exactly once"
    assert all(v == 42 for v in results)


def test_sequential_calls_still_once():
    calls = []
    lazy = Lazy(lambda: (calls.append(1), 7)[1])
    for _ in range(5):
        assert lazy.get() == 7
    assert len(calls) == 1
