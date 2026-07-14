def rotate_image(matrix: list[list[int]]) -> None:
    n = len(matrix)

    # Transpose the matrix (reflect across the main diagonal).
    for i in range(n):
        for j in range(i + 1, n):
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]

    # Reverse each row to complete the clockwise rotation.
    for row in matrix:
        row.reverse()
