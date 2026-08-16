from solution import smooth


def test_example_2x2():
    assert smooth([[1, 1], [1, 3]]) == [[1, 1], [1, 1]]


def test_example_3x3_cross():
    grid = [[0, 9, 0], [9, 0, 9], [0, 9, 0]]
    assert smooth(grid) == [[4, 4, 4], [4, 4, 4], [4, 4, 4]]


def test_example_single_row():
    assert smooth([[3, 6, 9]]) == [[4, 6, 7]]


def test_1x1_unchanged():
    assert smooth([[7]]) == [[7]]


def test_uniform_non_square_unchanged():
    grid = [[5, 5, 5, 5], [5, 5, 5, 5]]
    assert smooth(grid) == [[5, 5, 5, 5], [5, 5, 5, 5]]


def test_single_column():
    assert smooth([[1], [2], [3]]) == [[1], [2], [2]]


def test_reads_original_not_partial_result():
    # In-place smoothing left-to-right would yield [[5, 1, 0]].
    assert smooth([[10, 0, 0]]) == [[5, 3, 0]]


def test_input_not_mutated_and_new_grid_returned():
    grid = [[1, 1], [1, 3]]
    result = smooth(grid)
    assert grid == [[1, 1], [1, 3]]
    assert result is not grid
