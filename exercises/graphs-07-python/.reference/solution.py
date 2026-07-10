def solve(board: list[list[str]]) -> None:
    rows, cols = len(board), len(board[0])
    if rows == 0 or cols == 0:
        return

    def dfs(r: int, c: int) -> None:
        if r < 0 or r >= rows or c < 0 or c >= cols or board[r][c] != "O":
            return
        board[r][c] = "#"
        dfs(r + 1, c)
        dfs(r - 1, c)
        dfs(r, c + 1)
        dfs(r, c - 1)

    for c in range(cols):
        dfs(0, c)
        dfs(rows - 1, c)
    for r in range(rows):
        dfs(r, 0)
        dfs(r, cols - 1)

    for r in range(rows):
        for c in range(cols):
            if board[r][c] == "O":
                board[r][c] = "X"
            elif board[r][c] == "#":
                board[r][c] = "O"
