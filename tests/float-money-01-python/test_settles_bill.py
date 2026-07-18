from solution import settles_bill


def test_three_dimes_settle_thirty_cents():
    assert settles_bill([0.1, 0.1, 0.1], 0.3)


def test_exact_integer_amounts_settle():
    assert settles_bill([10.0, 20.0, 70.0], 100.0)


def test_settles_with_different_rounding_pattern():
    assert settles_bill([0.7, 0.1], 0.8)


def test_settles_another_rounding_pattern():
    assert settles_bill([1.1, 2.2], 3.3)


def test_short_by_a_cent_does_not_settle():
    assert not settles_bill([10.00, 10.00], 20.01)


def test_clearly_short_does_not_settle():
    assert not settles_bill([5.00, 5.00], 11.00)


def test_empty_amounts_settle_zero_bill():
    assert settles_bill([], 0.0)
