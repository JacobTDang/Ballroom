class DetectSquares:
    """Tracks added points and counts axis-aligned squares formable
    with a query point."""

    def __init__(self):
        self.points: dict[tuple[int, int], int] = {}

    def add(self, point: list[int]) -> None:
        raise NotImplementedError

    def count(self, point: list[int]) -> int:
        raise NotImplementedError
