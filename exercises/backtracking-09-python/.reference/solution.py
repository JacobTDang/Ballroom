def solve_n_queens(n: int) -> list[list[str]]:
    """Return every distinct board configuration that places n
    queens on an n x n board with no two attacking each other."""
    res: list[list[str]] = []
    cols: set[int] = set()
    diag1: set[int] = set()  # r+c
    diag2: set[int] = set()  # r-c
    placement: list[int] = []  # placement[r] = column of queen in row r

    def backtrack(r: int) -> None:
        if r == n:
            board = []
            for c in placement:
                row = "." * c + "Q" + "." * (n - c - 1)
                board.append(row)
            res.append(board)
            return
        for c in range(n):
            if c in cols or (r + c) in diag1 or (r - c) in diag2:
                continue
            cols.add(c)
            diag1.add(r + c)
            diag2.add(r - c)
            placement.append(c)
            backtrack(r + 1)
            placement.pop()
            cols.remove(c)
            diag1.remove(r + c)
            diag2.remove(r - c)

    backtrack(0)
    return res
