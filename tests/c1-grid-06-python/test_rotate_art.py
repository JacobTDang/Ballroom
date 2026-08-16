from solution import rotate


def test_example_2x2():
    assert rotate([[1, 2], [3, 4]]) == [[3, 1], [4, 2]]


def test_example_2x3_non_square():
    assert rotate([[1, 2, 3], [4, 5, 6]]) == [[4, 1], [5, 2], [6, 3]]


def test_example_single_row_becomes_column():
    assert rotate([[1, 2, 3]]) == [[1], [2], [3]]


def test_single_column_becomes_row():
    assert rotate([[1], [2], [3]]) == [[3, 2, 1]]


def test_1x1():
    assert rotate([[9]]) == [[9]]


def test_3x4_non_square():
    grid = [[1, 2, 3, 4], [5, 6, 7, 8], [9, 10, 11, 12]]
    assert rotate(grid) == [[9, 5, 1], [10, 6, 2], [11, 7, 3], [12, 8, 4]]


def test_four_rotations_identity():
    grid = [[1, 2, 3], [4, 5, 6]]
    assert rotate(rotate(rotate(rotate(grid)))) == [[1, 2, 3], [4, 5, 6]]


def test_input_not_mutated():
    grid = [[1, 2], [3, 4]]
    rotate(grid)
    assert grid == [[1, 2], [3, 4]]
