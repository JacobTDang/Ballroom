from solution import SlidingWindow


def test_boundary_burst_is_caught():
    s = SlidingWindow(2, 100)
    assert s.allow_at(90)
    assert s.allow_at(95)
    assert not s.allow_at(105), "burst straddling the old boundary must be denied"
    assert not s.allow_at(110)
    assert s.allow_at(191)


def test_requests_age_out_exactly():
    s = SlidingWindow(1, 100)
    assert s.allow_at(1000)
    assert not s.allow_at(1099), "99ms-old request still counts"
    assert s.allow_at(1100), "exactly-window-old request must no longer count"


def test_denied_requests_do_not_count():
    s = SlidingWindow(2, 100)
    s.allow_at(0)
    s.allow_at(1)
    for i in range(2, 50):
        assert not s.allow_at(i)
    assert s.allow_at(101), "denied requests must not consume the budget"


def test_steady_rate_under_limit_always_passes():
    s = SlidingWindow(2, 100)
    for at in range(0, 1000, 60):
        assert s.allow_at(at)
