from solution import min_eating_speed


def test_min_eating_speed_case_1():
    assert min_eating_speed([3, 6, 7, 11], 8) == 4


def test_min_eating_speed_case_2():
    assert min_eating_speed([30, 11, 23, 4, 20], 5) == 30


def test_min_eating_speed_case_3():
    assert min_eating_speed([30, 11, 23, 4, 20], 6) == 23


def test_min_eating_speed_case_4():
    assert min_eating_speed([1000000000], 2) == 500000000


def test_min_eating_speed_case_5():
    assert min_eating_speed([1], 1) == 1


def test_min_eating_speed_case_6():
    assert min_eating_speed([3, 6, 7, 11], 4) == 11


def test_min_eating_speed_case_7():
    assert min_eating_speed([5], 5) == 1


def test_min_eating_speed_case_8():
    assert min_eating_speed([1000000000], 1000000000) == 1
