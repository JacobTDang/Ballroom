import threading
import time

from solution import Server


def test_stop_drains_everything_accepted():
    handled = []
    lock = threading.Lock()

    def handle(v):
        time.sleep(0.002)
        with lock:
            handled.append(v)

    s = Server(4, handle)
    accepted = sum(1 for i in range(200) if s.submit(i))

    stopper = threading.Thread(target=s.stop)
    stopper.start()
    stopper.join(10)
    assert not stopper.is_alive(), "stop() never returned -- deadlock or never drained"
    assert len(handled) == accepted, f"{len(handled)} handled after stop, want every accepted job ({accepted})"


def test_submit_refused_after_stop():
    handled = []
    s = Server(2, handled.append)
    s.submit(1)
    s.stop()
    before = len(handled)

    assert not s.submit(2), "submit accepted a job after stop returned"
    time.sleep(0.02)
    assert len(handled) == before
