import threading
import time

from solution import Account, transfer


def test_crossed_transfers_do_not_deadlock():
    a = Account(1, 10000)
    b = Account(2, 10000)

    def worker(i):
        for _ in range(50):
            if i % 2 == 0:
                transfer(a, b, 1)
            else:
                transfer(b, a, 1)

    threads = [threading.Thread(target=worker, args=(i,), daemon=True) for i in range(4)]
    for t in threads:
        t.start()
    deadline = time.time() + 10
    for t in threads:
        t.join(max(0.1, deadline - time.time()))
    assert not any(t.is_alive() for t in threads), \
        "crossed transfers deadlocked (each direction holding one lock, waiting on the other)"

    assert a.balance() + b.balance() == 20000, "total balance not conserved"


def test_insufficient_funds_moves_nothing():
    a = Account(1, 5)
    b = Account(2, 0)
    assert not transfer(a, b, 10)
    assert a.balance() == 5 and b.balance() == 0
