from solution import jump


def test_classic():
    assert jump([2, 3, 1, 1, 4]) == 2


def test_single_element():
    assert jump([0]) == 0


def test_all_ones():
    assert jump([1, 1, 1, 1]) == 3


def test_big_first_jump():
    assert jump([5, 0, 0, 0, 0]) == 1


def test_two_elements():
    assert jump([2, 1]) == 1
