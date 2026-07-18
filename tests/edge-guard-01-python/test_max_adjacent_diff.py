import pytest

from solution import max_adjacent_diff


def test_max_adjacent_diff_case_1():
    assert max_adjacent_diff([3, 1, 4, 1, 5, 9, 2, 6]) == 7


def test_max_adjacent_diff_case_2():
    assert max_adjacent_diff([5, 5]) == 0


def test_max_adjacent_diff_case_3():
    assert max_adjacent_diff([-5, -1, -10]) == 9


def test_max_adjacent_diff_case_4():
    assert max_adjacent_diff([1, 100]) == 99


def test_max_adjacent_diff_empty_raises():
    with pytest.raises(ValueError):
        max_adjacent_diff([])


def test_max_adjacent_diff_single_element_raises():
    with pytest.raises(ValueError):
        max_adjacent_diff([42])
