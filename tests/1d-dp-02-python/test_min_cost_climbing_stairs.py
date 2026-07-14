from solution import min_cost_climbing_stairs


def test_three():
    assert min_cost_climbing_stairs([10, 15, 20]) == 15


def test_ten():
    assert min_cost_climbing_stairs([1, 100, 1, 1, 1, 100, 1, 1, 100, 1]) == 6


def test_two_equal():
    assert min_cost_climbing_stairs([0, 0]) == 0


def test_boundary_max_values():
    assert min_cost_climbing_stairs([999, 999]) == 999


def test_larger_ascending():
    assert min_cost_climbing_stairs([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]) == 25
