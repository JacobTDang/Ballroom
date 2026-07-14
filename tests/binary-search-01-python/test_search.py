from solution import search


def test_search():
    assert search([-1, 0, 3, 5, 9, 12], 9) == 4
    assert search([-1, 0, 3, 5, 9, 12], 2) == -1
    assert search([5], 5) == 0
    assert search([2, 5], 5) == 1
    assert search([2, 5], 1) == -1
    assert search([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 1) == 0
    assert search([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 10) == 9
    assert search([-10, -5, 0, 5, 10], -5) == 1
