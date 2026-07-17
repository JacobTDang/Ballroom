import threading
import time

from solution import process_all


def test_results_in_input_order():
    jobs = list(range(100))
    got = process_all(jobs, 8, lambda v: (time.sleep(0.002), v * 2)[1])
    assert got == [v * 2 for v in jobs]


def test_actually_parallel_within_bound():
    in_flight = 0
    high_water = 0
    lock = threading.Lock()

    def fn(v):
        nonlocal in_flight, high_water
        with lock:
            in_flight += 1
            high_water = max(high_water, in_flight)
        time.sleep(0.01)
        with lock:
            in_flight -= 1
        return v

    process_all(list(range(64)), 8, fn)

    assert high_water >= 2, f"high-water {high_water}: the pool never ran jobs in parallel"
    assert high_water <= 8, f"high-water {high_water}: more in flight than the 8 workers requested"
