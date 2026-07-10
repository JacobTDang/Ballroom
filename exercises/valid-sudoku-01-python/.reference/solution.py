def is_valid_sudoku(board: list[str]) -> bool:
    """Report whether the filled cells of a 9x9 Sudoku board satisfy
    Sudoku's placement rules (no digit repeated within a row, column, or
    3x3 box). Empty cells are '.'."""
    rows = [set() for _ in range(9)]
    cols = [set() for _ in range(9)]
    boxes = [set() for _ in range(9)]

    for r in range(9):
        for c in range(9):
            ch = board[r][c]
            if ch == ".":
                continue
            b = (r // 3) * 3 + c // 3
            if ch in rows[r] or ch in cols[c] or ch in boxes[b]:
                return False
            rows[r].add(ch)
            cols[c].add(ch)
            boxes[b].add(ch)
    return True
