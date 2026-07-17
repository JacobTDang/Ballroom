import threading
import time

from solution import run_limited


def make_instrumented(n, state):
    lock = threading.Lock()

    def task():
        with lock:
            state["in_flight"] += 1
            state["high_water"] = max(state["high_water"], state["in_flight"])
        time.sleep(0.015)
        with lock:
            state["in_flight"] -= 1
            state["ran"] += 1

    return [task for _ in range(n)]


def test_bound_holds_and_everything_runs():
    state = {"in_flight": 0, "high_water": 0, "ran": 0}
    run_limited(make_instrumented(32, state), 4)
    assert state["ran"] == 32
    assert state["high_water"] <= 4, f"high-water {state['high_water']} exceeded limit 4"
    assert state["high_water"] >= 2, f"high-water {state['high_water']}: no real parallelism"


def test_limit_one_is_serial():
    state = {"in_flight": 0, "high_water": 0, "ran": 0}
    run_limited(make_instrumented(6, state), 1)
    assert state["ran"] == 6
    assert state["high_water"] == 1


def test_limit_larger_than_tasks():
    state = {"in_flight": 0, "high_water": 0, "ran": 0}
    run_limited(make_instrumented(3, state), 10)
    assert state["ran"] == 3
