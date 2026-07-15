from solution import search_matrix

M = [[1, 3, 5, 7], [10, 11, 16, 20], [23, 30, 34, 60]]


def test_search_matrix_case_1():
    assert search_matrix(M, 3) is True


def test_search_matrix_case_2():
    assert search_matrix(M, 13) is False


def test_search_matrix_case_3():
    assert search_matrix([[1]], 1) is True


def test_search_matrix_case_4():
    assert search_matrix([[1, 3]], 3) is True


def test_search_matrix_case_5():
    assert search_matrix(M, 60) is True


def test_search_matrix_case_6():
    assert search_matrix(M, 0) is False


def test_search_matrix_case_7():
    assert search_matrix(M, 1) is True


def test_search_matrix_case_8():
    assert search_matrix(M, 23) is True
