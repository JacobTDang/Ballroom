from solution import is_n_straight_hand


def test_classic():
    assert is_n_straight_hand([1, 2, 3, 6, 2, 3, 4, 7, 8], 3) is True


def test_not_divisible():
    assert is_n_straight_hand([1, 2, 3, 4, 5], 4) is False


def test_missing_card():
    assert is_n_straight_hand([1, 2, 3, 4, 5, 7], 3) is False


def test_group_size_one():
    assert is_n_straight_hand([5, 5, 5], 1) is True


def test_exact_single_group():
    assert is_n_straight_hand([1, 2, 3], 3) is True
