from solution import search


def test_search_case_1():
    assert search([-1, 0, 3, 5, 9, 12], 9) == 4


def test_search_case_2():
    assert search([-1, 0, 3, 5, 9, 12], 2) == -1


def test_search_case_3():
    assert search([5], 5) == 0


def test_search_case_4():
    assert search([2, 5], 5) == 1


def test_search_case_5():
    assert search([2, 5], 1) == -1


def test_search_case_6():
    assert search([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 1) == 0


def test_search_case_7():
    assert search([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 10) == 9


def test_search_case_8():
    assert search([-10, -5, 0, 5, 10], -5) == 1
