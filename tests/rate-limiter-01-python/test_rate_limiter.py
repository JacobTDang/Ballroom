import time

from solution import RateLimiter


def test_allows_up_to_limit_within_window():
    rl = RateLimiter(limit=3, window_seconds=10)
    assert rl.allow() is True
    assert rl.allow() is True
    assert rl.allow() is True
    assert rl.allow() is False


def test_resets_after_window(monkeypatch):
    t = [1000.0]
    monkeypatch.setattr(time, "time", lambda: t[0])

    rl = RateLimiter(limit=1, window_seconds=5)
    assert rl.allow() is True
    assert rl.allow() is False

    t[0] += 5.1
    assert rl.allow() is True
