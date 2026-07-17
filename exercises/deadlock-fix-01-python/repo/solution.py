import threading
import time


class Account:
    def __init__(self, id: int, balance: int):
        self.id = id
        self._balance = balance
        self.lock = threading.Lock()

    def balance(self) -> int:
        with self.lock:
            return self._balance


def transfer(from_acct: Account, to_acct: Account, amount: int) -> bool:
    """Moves amount between accounts, locking both.

    TODO: locking from-then-to deadlocks the moment two transfers
    cross (A->B and B->A each hold one lock and wait for the other).
    Fix the ordering -- don't just wrap everything in one global lock.
    """
    with from_acct.lock:
        time.sleep(0.001)  # bookkeeping -- widens the inversion window
        with to_acct.lock:
            if from_acct._balance < amount:
                return False
            from_acct._balance -= amount
            to_acct._balance += amount
            return True
