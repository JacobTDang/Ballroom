from solution import is_valid_sudoku

VALID_BOARD = [
    "53..7....",
    "6..195...",
    ".98....6.",
    "8...6...3",
    "4..8.3..1",
    "7...2...6",
    ".6....28.",
    "...419..5",
    "....8..79",
]
INVALID_COLUMN_BOARD = [
    "83..7....",
    "6..195...",
    ".98....6.",
    "8...6...3",
    "4..8.3..1",
    "7...2...6",
    ".6....28.",
    "...419..5",
    "....8..79",
]
INVALID_ROW_BOARD = ["5.......5"] + ["........."] * 8
INVALID_BOX_BOARD = ["1........", ".1......."] + ["........."] * 7
EMPTY_BOARD = ["........."] * 9


def test_is_valid_sudoku():
    assert is_valid_sudoku(VALID_BOARD) is True
    assert is_valid_sudoku(INVALID_COLUMN_BOARD) is False
    assert is_valid_sudoku(INVALID_ROW_BOARD) is False
    assert is_valid_sudoku(INVALID_BOX_BOARD) is False
    assert is_valid_sudoku(EMPTY_BOARD) is True
