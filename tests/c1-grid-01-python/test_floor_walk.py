from solution import walk


def test_example_open_path():
    grid = [[0, 0, 1], [0, 1, 0], [0, 0, 0]]
    assert walk(grid, "DDRR") == (2, 2)


def test_example_wall_blocks():
    assert walk([[0, 1], [0, 0]], "RDR") == (1, 1)


def test_example_off_grid_ignored_non_square():
    assert walk([[0, 0]], "ULDRR") == (0, 1)


def test_1x1_all_moves_ignored():
    assert walk([[0]], "UDLRUDLR") == (0, 0)


def test_empty_moves_no_op():
    assert walk([[0, 0], [0, 0]], "") == (0, 0)


def test_walled_in_start_never_moves():
    assert walk([[0, 1], [1, 1]], "RDLURDLU") == (0, 0)


def test_single_column():
    grid = [[0], [0], [0]]
    assert walk(grid, "DDDDU") == (1, 0)


def test_blocked_move_then_continue():
    grid = [[0, 1, 0], [0, 0, 0]]
    assert walk(grid, "RDRRU") == (0, 2)


def test_non_square_tall_grid():
    grid = [[0, 0], [1, 0], [0, 0], [0, 1]]
    assert walk(grid, "RDDLDL") == (3, 0)
