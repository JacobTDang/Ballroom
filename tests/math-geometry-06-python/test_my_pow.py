import math

from solution import my_pow


def approx_equal(a: float, b: float) -> bool:
    return math.isclose(a, b, abs_tol=1e-6)


def test_my_pow_positive_exponent():
    assert approx_equal(my_pow(2.0, 10), 1024.0)


def test_my_pow_fractional_base():
    assert approx_equal(my_pow(2.1, 3), 9.261)


def test_my_pow_negative_exponent():
    assert approx_equal(my_pow(2.0, -2), 0.25)


def test_my_pow_zero_exponent():
    assert approx_equal(my_pow(0.5, 0), 1.0)


def test_my_pow_negative_base():
    assert approx_equal(my_pow(-2.0, 3), -8.0)


def test_my_pow_min_int32_exponent():
    # x = 1 keeps the expected result exact regardless of exponent
    # magnitude, while still exercising the negate-the-most-negative-
    # exponent overflow edge case.
    assert approx_equal(my_pow(1.0, -(2**31)), 1.0)


def test_my_pow_negative_base_negative_exponent():
    assert approx_equal(my_pow(-2.0, -2), 0.25)


def test_my_pow_larger_positive_exponent():
    assert approx_equal(my_pow(3.0, 5), 243.0)


def test_my_pow_fractional_squared():
    assert approx_equal(my_pow(1.5, 2), 2.25)
