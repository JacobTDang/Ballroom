from solution import max_of


def test_max_of_case_1():
    assert max_of([3, 1, 4, 1, 5, 9, 2, 6]) == 9


def test_max_of_case_2():
    assert max_of([-5, -1, -10]) == -1


def test_max_of_case_3():
    assert max_of([42]) == 42


def test_max_of_case_4():
    assert max_of([5, 5, 5]) == 5


def test_max_of_case_5():
    assert max_of([1, 2, 3, 4, 5, 100]) == 100


def test_max_of_case_6():
    assert max_of([-1, -2, -3, -100]) == -1
