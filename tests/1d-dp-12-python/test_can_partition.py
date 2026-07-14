from solution import can_partition


def test_classic():
    assert can_partition([1, 5, 11, 5]) is True


def test_odd_sum():
    assert can_partition([1, 2, 3, 5]) is False


def test_even_split():
    assert can_partition([1, 2, 3, 4]) is True


def test_two_equal():
    assert can_partition([2, 2]) is True


def test_single_element():
    assert can_partition([4]) is False


def test_all_same():
    assert can_partition([3, 3, 3, 3]) is True


def test_even_sum_unreachable():
    assert can_partition([2, 2, 3, 5]) is False


def test_boundary_values():
    assert can_partition([100, 100, 100, 100]) is True


def test_larger_multi_combination():
    assert can_partition([1, 2, 3, 4, 5, 6, 7]) is True
