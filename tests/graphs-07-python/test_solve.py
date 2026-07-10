from solution import solve


def to_grid(rows):
    return [list(row) for row in rows]


def test_classic():
    board = to_grid(["XXXX", "XOOX", "XXOX", "XOXX"])
    want = to_grid(["XXXX", "XXXX", "XXXX", "XOXX"])
    solve(board)
    assert board == want


def test_all_border_connected():
    board = to_grid(["OOO", "OXO", "OOO"])
    want = to_grid(["OOO", "OXO", "OOO"])
    solve(board)
    assert board == want


def test_single_cell():
    board = to_grid(["O"])
    want = to_grid(["O"])
    solve(board)
    assert board == want
