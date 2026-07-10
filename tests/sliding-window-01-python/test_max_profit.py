from solution import max_profit


def test_max_profit():
    assert max_profit([7, 1, 5, 3, 6, 4]) == 5
    assert max_profit([7, 6, 4, 3, 1]) == 0
    assert max_profit([2, 4, 1]) == 2
    assert max_profit([1]) == 0
    assert max_profit([3, 3, 3, 3]) == 0
    assert max_profit([1, 2, 4, 2, 5, 7, 2, 4, 9, 0]) == 8
