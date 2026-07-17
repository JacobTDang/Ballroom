import threading

from solution import TokenBucket


def hammer(bucket, callers):
    allowed = []
    lock = threading.Lock()
    # A start gate: every thread is created and waiting before any
    # calls allow(), so the calls genuinely overlap -- without it,
    # thread-startup serialization can hide the race entirely.
    gate = threading.Event()

    def caller():
        gate.wait()
        if bucket.allow():
            with lock:
                allowed.append(1)

    threads = [threading.Thread(target=caller) for _ in range(callers)]
    for t in threads:
        t.start()
    gate.set()
    for t in threads:
        t.join()
    return len(allowed)


def test_exactly_capacity_allowed_under_contention():
    b = TokenBucket(100)
    got = hammer(b, 300)
    assert got == 100, f"{got} of 300 concurrent allow() calls succeeded, want exactly 100"


def test_refill_grants_exactly_that_many():
    b = TokenBucket(100)
    hammer(b, 300)
    b.refill(40)
    got = hammer(b, 200)
    assert got == 40, f"{got} allowed after refill(40), want exactly 40"


def test_refill_clamps_at_capacity():
    b = TokenBucket(50)
    b.refill(1000)
    got = hammer(b, 200)
    assert got == 50, f"{got} allowed after over-refill, want capacity 50"
