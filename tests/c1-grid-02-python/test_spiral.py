from solution import spiral


def test_example_3x3():
    grid = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
    assert spiral(grid) == [1, 2, 3, 6, 9, 8, 7, 4, 5]


def test_example_2x3_non_square():
    assert spiral([[1, 2, 3], [4, 5, 6]]) == [1, 2, 3, 6, 5, 4]


def test_example_single_column():
    assert spiral([[7], [8], [9]]) == [7, 8, 9]


def test_1x1():
    assert spiral([[5]]) == [5]


def test_single_row():
    assert spiral([[1, 2, 3, 4]]) == [1, 2, 3, 4]


def test_2x2():
    assert spiral([[1, 2], [4, 3]]) == [1, 2, 3, 4]


def test_3x2_non_square_tall():
    assert spiral([[1, 2], [3, 4], [5, 6]]) == [1, 2, 4, 6, 5, 3]


def test_4x4_two_full_rings():
    grid = [
        [1, 2, 3, 4],
        [5, 6, 7, 8],
        [9, 10, 11, 12],
        [13, 14, 15, 16],
    ]
    assert spiral(grid) == [1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10]
