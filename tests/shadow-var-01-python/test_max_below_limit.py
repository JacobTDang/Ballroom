from solution import max_below_limit


def test_max_below_limit_exact_match():
    assert max_below_limit([3, 7, 2, 9, 5], 7) == 7


def test_max_below_limit_no_exact_match():
    assert max_below_limit([3, 7, 2, 9, 5], 6) == 5


def test_max_below_limit_none_qualify():
    assert max_below_limit([10, 20, 30], 5) == -1


def test_max_below_limit_negatives():
    assert max_below_limit([-5, -1, -10], -2) == -5


def test_max_below_limit_single_qualifying():
    assert max_below_limit([5], 10) == 5


def test_max_below_limit_single_non_qualifying():
    assert max_below_limit([15], 10) == -1
