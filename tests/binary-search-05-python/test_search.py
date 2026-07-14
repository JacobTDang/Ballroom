from solution import search


def test_search():
    assert search([4, 5, 6, 7, 0, 1, 2], 0) == 4
    assert search([4, 5, 6, 7, 0, 1, 2], 3) == -1
    assert search([1], 0) == -1
    assert search([5, 1, 3], 5) == 0
    assert search([1, 3], 3) == 1
    assert search([9, 10, 1, 2, 3, 4, 5, 6, 7, 8], 8) == 9
    assert search([9, 10, 1, 2, 3, 4, 5, 6, 7, 8], 100) == -1
