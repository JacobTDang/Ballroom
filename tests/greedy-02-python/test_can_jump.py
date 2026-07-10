from solution import can_jump


def test_classic():
    assert can_jump([2, 3, 1, 1, 4]) is True


def test_classic_false():
    assert can_jump([3, 2, 1, 0, 4]) is False


def test_single_element():
    assert can_jump([0]) is True


def test_zero_at_start_blocks():
    assert can_jump([0, 1]) is False


def test_big_first_jump_covers_rest():
    assert can_jump([5, 0, 0, 0, 0]) is True
