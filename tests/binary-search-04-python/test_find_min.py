from solution import find_min


def test_find_min():
    assert find_min([3, 4, 5, 1, 2]) == 1
    assert find_min([4, 5, 6, 7, 0, 1, 2]) == 0
    assert find_min([11, 13, 15, 17]) == 11
    assert find_min([2, 1]) == 1
    assert find_min([1]) == 1
    assert find_min([1, 2, 3, 4, 5]) == 1
    assert find_min([15, 18, 2, 3, 6, 12]) == 2
