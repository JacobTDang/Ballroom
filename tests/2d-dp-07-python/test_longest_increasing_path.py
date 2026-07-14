from solution import longest_increasing_path


def test_classic():
    assert longest_increasing_path([[9, 9, 4], [6, 6, 8], [2, 1, 1]]) == 4


def test_second_classic():
    assert longest_increasing_path([[3, 4, 5], [3, 2, 6], [2, 2, 1]]) == 4


def test_single_cell():
    assert longest_increasing_path([[1]]) == 1


def test_single_row():
    assert longest_increasing_path([[1, 2, 3, 4]]) == 4


def test_all_equal():
    assert longest_increasing_path([[7, 7], [7, 7]]) == 1


def test_single_column():
    assert longest_increasing_path([[1], [2], [3]]) == 3


def test_snake_full_traversal():
    assert longest_increasing_path([[1, 2, 3], [6, 5, 4], [7, 8, 9]]) == 9


def test_negative_values():
    assert longest_increasing_path([[-1, -2], [-3, -4]]) == 3
