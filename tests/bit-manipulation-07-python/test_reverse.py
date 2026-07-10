from solution import reverse


def test_classic():
    assert reverse(123) == 321


def test_negative():
    assert reverse(-123) == -321


def test_trailing_zero():
    assert reverse(120) == 21


def test_overflow_positive():
    assert reverse(1534236469) == 0


def test_overflow_negative():
    assert reverse(-2147483648) == 0


def test_zero():
    assert reverse(0) == 0
