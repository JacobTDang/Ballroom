from solution import find_redundant_connection


def test_triangle():
    assert find_redundant_connection([[1, 2], [1, 3], [2, 3]]) == [2, 3]


def test_later_cycle():
    edges = [[1, 2], [2, 3], [3, 4], [1, 4], [1, 5]]
    assert find_redundant_connection(edges) == [1, 4]


def test_merging_components():
    edges = [[1, 4], [3, 4], [1, 3], [1, 2], [4, 5]]
    assert find_redundant_connection(edges) == [1, 3]
