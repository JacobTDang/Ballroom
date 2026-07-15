from solution import find_min


def test_find_min_case_1():
    assert find_min([3, 4, 5, 1, 2]) == 1


def test_find_min_case_2():
    assert find_min([4, 5, 6, 7, 0, 1, 2]) == 0


def test_find_min_case_3():
    assert find_min([11, 13, 15, 17]) == 11


def test_find_min_case_4():
    assert find_min([2, 1]) == 1


def test_find_min_case_5():
    assert find_min([1]) == 1


def test_find_min_case_6():
    assert find_min([1, 2, 3, 4, 5]) == 1


def test_find_min_case_7():
    assert find_min([15, 18, 2, 3, 6, 12]) == 2
