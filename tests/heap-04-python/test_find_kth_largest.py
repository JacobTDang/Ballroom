from solution import find_kth_largest


def test_find_kth_largest():
    assert find_kth_largest([3, 2, 1, 5, 6, 4], 2) == 5
    assert find_kth_largest([3, 2, 3, 1, 2, 4, 5, 5, 6], 4) == 4
    assert find_kth_largest([1], 1) == 1
    assert find_kth_largest([2, 1], 2) == 1
    assert find_kth_largest([5, 5, 5, 5], 2) == 5
    assert find_kth_largest([-1, -5, -3, -2, -4], 3) == -3
