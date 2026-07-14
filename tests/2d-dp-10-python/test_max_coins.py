from solution import max_coins


def test_classic():
    assert max_coins([3, 1, 5, 8]) == 167


def test_two_balloons():
    assert max_coins([1, 5]) == 10


def test_single_balloon():
    assert max_coins([7]) == 7


def test_ones():
    assert max_coins([1, 1]) == 2


def test_all_ones_larger():
    assert max_coins([1, 1, 1, 1]) == 4


def test_zero_value_balloon():
    assert max_coins([3, 0, 5]) == 20


def test_larger_ascending():
    assert max_coins([1, 2, 3, 4, 5]) == 110


def test_descending():
    assert max_coins([5, 3, 1]) == 25
