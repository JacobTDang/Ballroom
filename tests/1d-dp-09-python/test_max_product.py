from solution import max_product


def test_classic():
    assert max_product([2, 3, -2, 4]) == 6


def test_zero_splits():
    assert max_product([-2, 0, -1]) == 0


def test_two_negatives_flip():
    assert max_product([-2, 3, -4]) == 24


def test_single_negative():
    assert max_product([-5]) == -5


def test_all_positive():
    assert max_product([1, 2, 3, 4]) == 24


def test_single_positive():
    assert max_product([7]) == 7


def test_multiple_zero_split_islands():
    assert max_product([0, 2, 0, 3, 0]) == 3


def test_whole_array_even_negatives():
    assert max_product([2, -3, 4, -5]) == 120
