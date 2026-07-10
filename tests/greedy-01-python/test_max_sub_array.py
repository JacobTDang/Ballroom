from solution import max_sub_array


def test_classic():
    assert max_sub_array([-2, 1, -3, 4, -1, 2, 1, -5, 4]) == 6


def test_all_negative():
    assert max_sub_array([-3, -2, -1]) == -1


def test_all_positive():
    assert max_sub_array([1, 2, 3, 4]) == 10


def test_single_element():
    assert max_sub_array([5]) == 5


def test_large_negative_in_middle():
    assert max_sub_array([5, 4, -20, 7, 8]) == 15
