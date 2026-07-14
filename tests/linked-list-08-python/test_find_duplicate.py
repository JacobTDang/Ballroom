from solution import find_duplicate


def test_find_duplicate():
    assert find_duplicate([1, 3, 4, 2, 2]) == 2
    assert find_duplicate([3, 1, 3, 4, 2]) == 3
    assert find_duplicate([1, 1]) == 1
    assert find_duplicate([1, 1, 2]) == 1
    assert find_duplicate([2, 2, 2, 2, 2]) == 2
    assert find_duplicate([1, 2, 3, 4, 5, 6, 7, 8, 9, 5]) == 5
    assert find_duplicate([1, 2, 2]) == 2
