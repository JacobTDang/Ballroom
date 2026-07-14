from solution import max_profit


def test_classic():
    assert max_profit([1, 2, 3, 0, 2]) == 3


def test_single_day():
    assert max_profit([1]) == 0


def test_monotonic_increasing():
    assert max_profit([1, 2, 4]) == 3


def test_empty():
    assert max_profit([]) == 0


def test_monotonic_decreasing():
    assert max_profit([5, 4, 3, 2, 1]) == 0


def test_two_days_profit():
    assert max_profit([1, 2]) == 1


def test_cooldown_forces_wait():
    assert max_profit([1, 4, 2, 7]) == 6


def test_larger_multi_trade():
    assert max_profit([6, 1, 3, 2, 4, 7]) == 6


def test_boundary_values():
    assert max_profit([10000, 1]) == 0
