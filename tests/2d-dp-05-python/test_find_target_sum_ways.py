from solution import find_target_sum_ways


def test_classic():
    assert find_target_sum_ways([1, 1, 1, 1, 1], 3) == 5


def test_single():
    assert find_target_sum_ways([1], 1) == 1


def test_unreachable():
    assert find_target_sum_ways([1, 2, 3], 100) == 0


def test_zeros():
    assert find_target_sum_ways([0, 0, 0, 0, 0, 0, 0, 0, 1], 1) == 256


def test_zero_target():
    assert find_target_sum_ways([1, 1], 0) == 2


def test_negative_target():
    assert find_target_sum_ways([1, 1, 1, 1, 1], -3) == 5


def test_mixed_zero_target():
    assert find_target_sum_ways([1, 2, 1], 0) == 2


def test_parity_impossible():
    assert find_target_sum_ways([1, 2, 3], 1) == 0
