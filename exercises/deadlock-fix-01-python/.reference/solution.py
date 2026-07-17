import threading


class Account:
    def __init__(self, id: int, balance: int):
        self.id = id
        self._balance = balance
        self.lock = threading.Lock()

    def balance(self) -> int:
        with self.lock:
            return self._balance


def transfer(from_acct: Account, to_acct: Account, amount: int) -> bool:
    """Deadlock needs a cycle; acquiring locks in a fixed global order
    (by account id) makes cycles impossible -- crossed transfers take
    the same lock first and simply queue. Disjoint transfers still run
    fully in parallel."""
    first, second = sorted((from_acct, to_acct), key=lambda a: a.id)
    with first.lock:
        with second.lock:
            if from_acct._balance < amount:
                return False
            from_acct._balance -= amount
            to_acct._balance += amount
            return True
