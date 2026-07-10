from solution import search_matrix

M = [[1, 3, 5, 7], [10, 11, 16, 20], [23, 30, 34, 60]]


def test_search_matrix():
    assert search_matrix(M, 3) is True
    assert search_matrix(M, 13) is False
    assert search_matrix([[1]], 1) is True
    assert search_matrix([[1, 3]], 3) is True
    assert search_matrix(M, 60) is True
    assert search_matrix(M, 0) is False
