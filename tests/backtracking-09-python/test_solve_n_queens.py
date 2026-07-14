from solution import solve_n_queens


def normalize_exact(boards):
    return sorted(boards)


def test_solve_n_queens_four():
    got = normalize_exact(solve_n_queens(4))
    want = normalize_exact(
        [
            [".Q..", "...Q", "Q...", "..Q."],
            ["..Q.", "Q...", "...Q", ".Q.."],
        ]
    )
    assert got == want


def test_solve_n_queens_one():
    assert solve_n_queens(1) == [["Q"]]


def test_solve_n_queens_no_solutions_for_two_or_three():
    assert solve_n_queens(2) == []
    assert solve_n_queens(3) == []


def test_solve_n_queens_five_has_ten_solutions():
    got = solve_n_queens(5)
    assert len(got) == 10
    for board in got:
        assert len(board) == 5
        for row in board:
            assert len(row) == 5
