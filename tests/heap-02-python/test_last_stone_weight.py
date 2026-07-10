from solution import last_stone_weight


def test_last_stone_weight():
    assert last_stone_weight([2, 7, 4, 1, 8, 1]) == 1
    assert last_stone_weight([1]) == 1
    assert last_stone_weight([1, 1]) == 0
    assert last_stone_weight([1, 3]) == 2
    assert last_stone_weight([2, 2]) == 0
