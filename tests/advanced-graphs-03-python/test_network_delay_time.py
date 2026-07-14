from solution import network_delay_time


def test_classic():
    times = [[2, 1, 1], [2, 3, 1], [3, 4, 1]]
    assert network_delay_time(times, 4, 2) == 2


def test_single_edge_reachable():
    assert network_delay_time([[1, 2, 1]], 2, 1) == 1


def test_unreachable():
    assert network_delay_time([[1, 2, 1]], 2, 2) == -1


def test_shortest_of_multiple_paths():
    times = [[1, 2, 1], [2, 3, 2], [1, 3, 4]]
    assert network_delay_time(times, 3, 1) == 3


def test_single_node_no_edges():
    assert network_delay_time([], 1, 1) == 0
