def pacific_atlantic(heights: list[list[int]]) -> list[list[int]]:
    rows, cols = len(heights), len(heights[0])
    pacific = [[False] * cols for _ in range(rows)]
    atlantic = [[False] * cols for _ in range(rows)]

    def dfs(r: int, c: int, visited: list[list[bool]], prev_height: int) -> None:
        if (
            r < 0
            or r >= rows
            or c < 0
            or c >= cols
            or visited[r][c]
            or heights[r][c] < prev_height
        ):
            return
        visited[r][c] = True
        dfs(r + 1, c, visited, heights[r][c])
        dfs(r - 1, c, visited, heights[r][c])
        dfs(r, c + 1, visited, heights[r][c])
        dfs(r, c - 1, visited, heights[r][c])

    for c in range(cols):
        dfs(0, c, pacific, heights[0][c])
        dfs(rows - 1, c, atlantic, heights[rows - 1][c])
    for r in range(rows):
        dfs(r, 0, pacific, heights[r][0])
        dfs(r, cols - 1, atlantic, heights[r][cols - 1])

    res = []
    for r in range(rows):
        for c in range(cols):
            if pacific[r][c] and atlantic[r][c]:
                res.append([r, c])
    return res
