from solution import find_kth_largest


def test_find_kth_largest_case_1():
    assert find_kth_largest([3, 2, 1, 5, 6, 4], 2) == 5


def test_find_kth_largest_case_2():
    assert find_kth_largest([3, 2, 3, 1, 2, 4, 5, 5, 6], 4) == 4


def test_find_kth_largest_case_3():
    assert find_kth_largest([1], 1) == 1


def test_find_kth_largest_case_4():
    assert find_kth_largest([2, 1], 2) == 1


def test_find_kth_largest_case_5():
    assert find_kth_largest([5, 5, 5, 5], 2) == 5


def test_find_kth_largest_case_6():
    assert find_kth_largest([-1, -5, -3, -2, -4], 3) == -3
