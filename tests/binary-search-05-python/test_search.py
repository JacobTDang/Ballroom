from solution import search


def test_search_case_1():
    assert search([4, 5, 6, 7, 0, 1, 2], 0) == 4


def test_search_case_2():
    assert search([4, 5, 6, 7, 0, 1, 2], 3) == -1


def test_search_case_3():
    assert search([1], 0) == -1


def test_search_case_4():
    assert search([5, 1, 3], 5) == 0


def test_search_case_5():
    assert search([1, 3], 3) == 1


def test_search_case_6():
    assert search([9, 10, 1, 2, 3, 4, 5, 6, 7, 8], 8) == 9


def test_search_case_7():
    assert search([9, 10, 1, 2, 3, 4, 5, 6, 7, 8], 100) == -1
