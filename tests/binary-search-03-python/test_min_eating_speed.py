from solution import min_eating_speed


def test_min_eating_speed():
    assert min_eating_speed([3, 6, 7, 11], 8) == 4
    assert min_eating_speed([30, 11, 23, 4, 20], 5) == 30
    assert min_eating_speed([30, 11, 23, 4, 20], 6) == 23
    assert min_eating_speed([1000000000], 2) == 500000000
    assert min_eating_speed([1], 1) == 1
    assert min_eating_speed([3, 6, 7, 11], 4) == 11
    assert min_eating_speed([5], 5) == 1
    assert min_eating_speed([1000000000], 1000000000) == 1
