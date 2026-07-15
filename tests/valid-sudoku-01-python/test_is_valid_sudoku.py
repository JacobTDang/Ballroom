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
SOLVED_BOARD = [
    "534678912",
    "672195348",
    "198342567",
    "859761423",
    "426853791",
    "713924856",
    "961537284",
    "287419635",
    "345286179",
]
SAME_DIGIT_DIFFERENT_UNITS_BOARD = ["5........"] + ["........."] * 3 + ["....5...."] + ["........."] * 4
SINGLE_CELL_BOARD = ["5........"] + ["........."] * 8


def test_is_valid_sudoku_case_1():
    assert is_valid_sudoku(VALID_BOARD) is True


def test_is_valid_sudoku_case_2():
    assert is_valid_sudoku(INVALID_COLUMN_BOARD) is False


def test_is_valid_sudoku_case_3():
    assert is_valid_sudoku(INVALID_ROW_BOARD) is False


def test_is_valid_sudoku_case_4():
    assert is_valid_sudoku(INVALID_BOX_BOARD) is False


def test_is_valid_sudoku_case_5():
    assert is_valid_sudoku(EMPTY_BOARD) is True


def test_is_valid_sudoku_case_6():
    assert is_valid_sudoku(SOLVED_BOARD) is True


def test_is_valid_sudoku_case_7():
    assert is_valid_sudoku(SAME_DIGIT_DIFFERENT_UNITS_BOARD) is True


def test_is_valid_sudoku_case_8():
    assert is_valid_sudoku(SINGLE_CELL_BOARD) is True
