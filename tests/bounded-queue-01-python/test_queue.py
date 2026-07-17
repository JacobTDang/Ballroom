import threading
import time

from solution import BoundedQueue


def test_every_item_arrives_exactly_once():
    q = BoundedQueue(4)
    producers, per_producer, consumers = 3, 200, 3
    total = producers * per_producer

    def produce(p):
        for i in range(per_producer):
            q.put(p * per_producer + i)

    seen = []
    seen_lock = threading.Lock()

    def consume(n):
        for _ in range(n):
            v = q.get()
            with seen_lock:
                seen.append(v)

    threads = [threading.Thread(target=produce, args=(p,)) for p in range(producers)]
    threads += [threading.Thread(target=consume, args=(total // consumers,)) for _ in range(consumers)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert sorted(seen) == list(range(total))


def test_get_blocks_until_put():
    q = BoundedQueue(2)
    got = []
    done = threading.Event()

    def getter():
        got.append(q.get())
        done.set()

    t = threading.Thread(target=getter)
    t.start()
    assert not done.wait(0.05), "get() on an empty queue returned immediately"
    q.put(7)
    assert done.wait(1.0), "get() never woke up after a put"
    assert got == [7]
    t.join()


def test_put_blocks_until_get():
    q = BoundedQueue(2)
    q.put(1)
    q.put(2)
    done = threading.Event()

    def putter():
        q.put(3)
        done.set()

    t = threading.Thread(target=putter)
    t.start()
    assert not done.wait(0.05), "put() on a full queue returned immediately"
    assert q.get() == 1, "expected FIFO order"
    assert done.wait(1.0), "blocked put() never completed after a get freed a slot"
    t.join()
