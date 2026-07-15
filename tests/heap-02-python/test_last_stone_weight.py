from solution import last_stone_weight


def test_last_stone_weight_case_1():
    assert last_stone_weight([2, 7, 4, 1, 8, 1]) == 1


def test_last_stone_weight_case_2():
    assert last_stone_weight([1]) == 1


def test_last_stone_weight_case_3():
    assert last_stone_weight([1, 1]) == 0


def test_last_stone_weight_case_4():
    assert last_stone_weight([1, 3]) == 2


def test_last_stone_weight_case_5():
    assert last_stone_weight([2, 2]) == 0


def test_last_stone_weight_case_6():
    assert last_stone_weight([10, 4, 2, 10]) == 2


def test_last_stone_weight_case_7():
    assert last_stone_weight([1, 1, 1, 1]) == 0
