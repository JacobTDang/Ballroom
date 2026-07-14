def set_zeroes(matrix: list[list[int]]) -> None:
    rows = len(matrix)
    if rows == 0:
        return
    cols = len(matrix[0])

    zero_row = [False] * rows
    zero_col = [False] * cols

    for r in range(rows):
        for c in range(cols):
            if matrix[r][c] == 0:
                zero_row[r] = True
                zero_col[c] = True

    for r in range(rows):
        for c in range(cols):
            if zero_row[r] or zero_col[c]:
                matrix[r][c] = 0
