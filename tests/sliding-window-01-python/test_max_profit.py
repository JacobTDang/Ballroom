from solution import max_profit


def test_max_profit_case_1():
    assert max_profit([7, 1, 5, 3, 6, 4]) == 5


def test_max_profit_case_2():
    assert max_profit([7, 6, 4, 3, 1]) == 0


def test_max_profit_case_3():
    assert max_profit([2, 4, 1]) == 2


def test_max_profit_case_4():
    assert max_profit([1]) == 0


def test_max_profit_case_5():
    assert max_profit([3, 3, 3, 3]) == 0


def test_max_profit_case_6():
    assert max_profit([1, 2, 4, 2, 5, 7, 2, 4, 9, 0]) == 8


def test_max_profit_case_7():
    assert max_profit([]) == 0


def test_max_profit_case_8():
    assert max_profit([3, 1, 4, 1, 5, 9, 2, 6]) == 8


def test_max_profit_case_9():
    assert max_profit([2, 1, 2, 1, 0, 1, 2]) == 2
