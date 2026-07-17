import threading
import time

from solution import Barrier


def test_no_one_proceeds_early_across_rounds():
    n, rounds = 4, 5
    b = Barrier(n)
    arrivals = [0] * rounds
    lock = threading.Lock()
    failures = []

    def participant(p):
        for r in range(rounds):
            time.sleep(p * 0.003)  # staggered arrivals widen the window
            with lock:
                arrivals[r] += 1
            b.wait()
            with lock:
                if arrivals[r] != n:
                    failures.append(f"round {r}: proceeded with {arrivals[r]}/{n} arrivals")
                    return

    threads = [threading.Thread(target=participant, args=(p,)) for p in range(n)]
    for t in threads:
        t.start()
    deadline = time.time() + 10
    for t in threads:
        t.join(max(0.1, deadline - time.time()))
    assert not any(t.is_alive() for t in threads), "barrier deadlocked"
    assert not failures, failures[0]
