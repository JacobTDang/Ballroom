class DetectSquares:
    """Tracks added points and counts axis-aligned squares formable
    with a query point, using a frequency map from point to the number
    of times it was added. For a query point, every previously added
    point sharing its x-coordinate forms a candidate vertical edge; the
    two horizontal partner corners at that same side length are then
    checked for existence."""

    def __init__(self):
        self.points: dict[tuple[int, int], int] = {}

    def add(self, point: list[int]) -> None:
        key = (point[0], point[1])
        self.points[key] = self.points.get(key, 0) + 1

    def count(self, point: list[int]) -> int:
        qx, qy = point[0], point[1]
        total = 0

        for (px, py), freq in self.points.items():
            if px != qx or py == qy:
                continue
            side = py - qy
            for cx in (qx + side, qx - side):
                corner1 = (cx, qy)
                corner2 = (cx, py)
                if corner1 in self.points and corner2 in self.points:
                    total += freq * self.points[corner1] * self.points[corner2]

        return total
