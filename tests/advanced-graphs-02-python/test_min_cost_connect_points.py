from solution import min_cost_connect_points


def test_classic():
    points = [[0, 0], [2, 2], [3, 10], [5, 2], [7, 0]]
    assert min_cost_connect_points(points) == 20


def test_three_points():
    points = [[3, 12], [-2, 5], [-4, 1]]
    assert min_cost_connect_points(points) == 18


def test_single_point():
    assert min_cost_connect_points([[0, 0]]) == 0
