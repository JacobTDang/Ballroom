from solution import count_regions


def test_example_two_regions():
    assert count_regions([[1, 1, 0], [0, 1, 0], [0, 0, 1]]) == 2


def test_example_diagonals_do_not_connect():
    assert count_regions([[1, 0, 1], [0, 1, 0], [1, 0, 1]]) == 5


def test_example_single_row():
    assert count_regions([[1, 1, 1, 1]]) == 1


def test_single_row_with_gap():
    assert count_regions([[1, 0, 1, 1]]) == 2


def test_1x1():
    assert count_regions([[0]]) == 0
    assert count_regions([[1]]) == 1


def test_all_zeros_non_square():
    assert count_regions([[0, 0, 0], [0, 0, 0]]) == 0


def test_all_ones_non_square():
    assert count_regions([[1, 1], [1, 1], [1, 1]]) == 1


def test_single_column():
    assert count_regions([[1], [0], [1]]) == 2


def test_long_snake_single_region():
    size = 100
    grid = [[0] * size for _ in range(size)]
    for r in range(0, size, 2):
        for c in range(size):
            grid[r][c] = 1
    for r in range(1, size, 2):
        grid[r][size - 1 if (r // 2) % 2 == 0 else 0] = 1
    assert count_regions(grid) == 1
