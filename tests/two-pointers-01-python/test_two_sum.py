from solution import two_sum


def test_two_sum_case_1():
    assert two_sum([2, 7, 11, 15], 9) == [1, 2]


def test_two_sum_case_2():
    assert two_sum([2, 3, 4], 6) == [1, 3]


def test_two_sum_case_3():
    assert two_sum([-1, 0], -1) == [1, 2]


def test_two_sum_case_4():
    assert two_sum([3, 3], 6) == [1, 2]


def test_two_sum_case_5():
    assert two_sum([1, 2, 3, 4, 4, 9, 56, 90], 8) == [4, 5]


def test_two_sum_case_6():
    assert two_sum([-8, -3, 0, 4, 9, 13], 5) == [1, 6]


def test_two_sum_case_7():
    assert two_sum([-8, -3, 0, 4, 9, 13], 4) == [3, 4]


def test_two_sum_case_8():
    assert two_sum([-3, -1, 0, 2, 4, 5], 1) == [1, 5]


def test_two_sum_case_9():
    assert two_sum([5, 25, 75], 100) == [2, 3]
