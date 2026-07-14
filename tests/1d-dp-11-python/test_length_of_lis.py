from solution import length_of_lis


def test_classic():
    assert length_of_lis([10, 9, 2, 5, 3, 7, 101, 18]) == 4


def test_repeated_dip():
    assert length_of_lis([0, 1, 0, 3, 2, 3]) == 4


def test_all_equal():
    assert length_of_lis([7, 7, 7, 7]) == 1


def test_single_element():
    assert length_of_lis([5]) == 1


def test_strictly_decreasing():
    assert length_of_lis([5, 4, 3, 2, 1]) == 1


def test_strictly_increasing():
    assert length_of_lis([1, 2, 3, 4, 5]) == 5


def test_negative_values():
    assert length_of_lis([-1, -2, 0, 1, -3, 5]) == 4


def test_boundary_values():
    assert length_of_lis([-10000, 10000]) == 2


def test_duplicates_break_streak():
    assert length_of_lis([3, 3, 3, 4, 4, 5]) == 3
