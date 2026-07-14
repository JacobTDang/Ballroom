from solution import coin_change


def test_classic():
    assert coin_change([1, 2, 5], 11) == 3


def test_impossible():
    assert coin_change([2], 3) == -1


def test_zero_amount():
    assert coin_change([1], 0) == 0


def test_single_coin_exact():
    assert coin_change([3, 7], 6) == 2


def test_large_amount_only_ones():
    assert coin_change([1], 10000) == 10000


def test_unreachable_amount():
    assert coin_change([3, 5], 7) == -1


def test_mixed_denominations():
    assert coin_change([1, 5, 10, 25], 63) == 6
