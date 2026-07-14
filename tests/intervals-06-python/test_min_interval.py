from solution import min_interval


def test_classic():
    assert min_interval([[1, 4], [2, 4], [3, 6], [4, 4]], [2, 3, 4, 5]) == [3, 3, 1, 4]


def test_some_queries_uncovered():
    assert min_interval([[2, 3], [2, 5], [1, 8], [20, 25]], [2, 19, 5, 22]) == [2, -1, 4, 6]


def test_boundary_endpoints_and_a_miss():
    assert min_interval([[1, 10]], [1, 10, 11]) == [10, 10, -1]


def test_single_point_interval():
    assert min_interval([[5, 5]], [5]) == [1]


def test_boundary_constraint_values():
    assert min_interval([[1, 10000000]], [1, 10000000]) == [10000000, 10000000]


def test_nested_intervals_varying_size_multiple_queries():
    intervals = [[1, 100], [10, 20], [15, 16], [50, 60]]
    assert min_interval(intervals, [15, 55, 99, 5]) == [2, 11, 100, 100]
