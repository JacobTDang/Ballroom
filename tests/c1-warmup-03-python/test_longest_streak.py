from solution import longest_streak


def test_example_streak_at_start():
    assert longest_streak([10, 20, 30, 5, 6]) == 3


def test_example_all_equal():
    assert longest_streak([7, 7, 7]) == 1


def test_example_empty():
    assert longest_streak([]) == 0


def test_single_element():
    assert longest_streak([42]) == 1


def test_all_increasing():
    assert longest_streak([1, 2, 3, 4, 5]) == 5


def test_all_decreasing():
    assert longest_streak([5, 4, 3, 2, 1]) == 1


def test_streak_at_end():
    assert longest_streak([9, 8, 1, 2, 3, 4]) == 4


def test_equal_breaks_streak():
    assert longest_streak([1, 2, 2, 3]) == 2
