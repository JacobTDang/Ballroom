from solution import remove_value


def test_remove_value_leading_run():
    assert remove_value([1, 2, 2, 2, 3], 2) == [1, 3]


def test_remove_value_split_run():
    assert remove_value([2, 2, 5, 2], 2) == [5]


def test_remove_value_no_match():
    assert remove_value([1, 3, 5], 9) == [1, 3, 5]


def test_remove_value_all_match():
    assert remove_value([2, 2, 2, 2], 2) == []


def test_remove_value_matches_in_middle():
    assert remove_value([5, 2, 2, 5], 2) == [5, 5]


def test_remove_value_single_match():
    assert remove_value([2], 2) == []


def test_remove_value_empty():
    assert remove_value([], 2) == []
