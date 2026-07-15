from solution import find_duplicate


def test_find_duplicate_case_1():
    assert find_duplicate([1, 3, 4, 2, 2]) == 2


def test_find_duplicate_case_2():
    assert find_duplicate([3, 1, 3, 4, 2]) == 3


def test_find_duplicate_case_3():
    assert find_duplicate([1, 1]) == 1


def test_find_duplicate_case_4():
    assert find_duplicate([1, 1, 2]) == 1


def test_find_duplicate_case_5():
    assert find_duplicate([2, 2, 2, 2, 2]) == 2


def test_find_duplicate_case_6():
    assert find_duplicate([1, 2, 3, 4, 5, 6, 7, 8, 9, 5]) == 5


def test_find_duplicate_case_7():
    assert find_duplicate([1, 2, 2]) == 2
