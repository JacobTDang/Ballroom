from solution import step


def test_example_simple_move():
    assert step([[">", ".", "."]], 1) == [[".", ">", "."]]


def test_example_blocked_pair():
    assert step([[">", ">", "."]], 1) == [[">", ".", ">"]]


def test_example_exits_at_right_edge():
    grid = [[">", ">"], [".", ">"]]
    assert step(grid, 1) == [[">", "."], [".", "."]]


def test_blocked_even_though_blocker_exits():
    assert step([[">", ">"]], 1) == [[">", "."]]


def test_zero_ticks_no_op_and_input_untouched():
    grid = [[">", "."], [".", ">"]]
    result = step(grid, 0)
    assert result == [[">", "."], [".", ">"]]
    assert grid == [[">", "."], [".", ">"]]


def test_1x1_marker_exits():
    assert step([[">"]], 1) == [["."]]


def test_all_empty_non_square_stays_empty():
    assert step([[".", ".", "."], [".", ".", "."]], 4) == [
        [".", ".", "."],
        [".", ".", "."],
    ]


def test_single_column_all_exit():
    assert step([[">"], ["."], [">"]], 1) == [["."], ["."], ["."]]


def test_full_row_drains_over_ticks():
    assert step([[">", ">", ">"]], 3) == [[".", ">", "."]]
    assert step([[">", ">", ">"]], 5) == [[".", ".", "."]]


def test_rows_move_independently():
    grid = [[">", "."], [">", ">"]]
    assert step(grid, 1) == [[".", ">"], [">", "."]]
