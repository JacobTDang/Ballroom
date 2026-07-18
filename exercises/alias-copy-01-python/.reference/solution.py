class Grid:
    def __init__(self, rows: int, cols: int):
        self.rows = [[0] * cols for _ in range(rows)]

    def get(self, r: int, c: int) -> int:
        return self.rows[r][c]

    def set(self, r: int, c: int, v: int) -> None:
        self.rows[r][c] = v

    def snapshot(self) -> list[list[int]]:
        """Returns an independent copy of the grid's current cell
        values, for a caller that wants to edit the live grid and
        later compare against this saved state."""
        return [row[:] for row in self.rows]
