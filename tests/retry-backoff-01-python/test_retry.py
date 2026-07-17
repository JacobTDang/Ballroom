import pytest

from solution import retry


def make_op(fail_times):
    state = {"calls": 0}

    def op():
        state["calls"] += 1
        if state["calls"] <= fail_times:
            raise RuntimeError("boom")
        return "ok"

    return op, state


def test_first_try_success_never_sleeps():
    slept = []
    op, state = make_op(0)
    assert retry(op, 5, 100, 1000, slept.append) == "ok"
    assert state["calls"] == 1
    assert slept == []


def test_exponential_delays_exact():
    slept = []
    op, state = make_op(3)
    assert retry(op, 5, 100, 10000, slept.append) == "ok"
    assert state["calls"] == 4
    assert slept == [100, 200, 400]


def test_cap_flattens_the_curve():
    slept = []
    op, _ = make_op(4)
    assert retry(op, 6, 100, 250, slept.append) == "ok"
    assert slept == [100, 200, 250, 250]


def test_exhaustion_reraises_last_error_no_trailing_sleep():
    slept = []
    calls = {"n": 0}

    def op():
        calls["n"] += 1
        raise RuntimeError("always")

    with pytest.raises(RuntimeError, match="always"):
        retry(op, 3, 100, 1000, slept.append)
    assert calls["n"] == 3
    assert slept == [100, 200], "never sleep after the final failure"
