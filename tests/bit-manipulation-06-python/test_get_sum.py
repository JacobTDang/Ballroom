from solution import get_sum


def test_classic():
    assert get_sum(1, 1) == 2


def test_positive_positive():
    assert get_sum(2, 3) == 5


def test_negative_positive_cancel():
    assert get_sum(-1, 1) == 0


def test_two_negatives():
    assert get_sum(-5, -7) == -12


def test_with_zero():
    assert get_sum(0, 0) == 0


def test_max_int32_bounds():
    assert get_sum(2147483647, 0) == 2147483647
    assert get_sum(-2147483648, 0) == -2147483648
