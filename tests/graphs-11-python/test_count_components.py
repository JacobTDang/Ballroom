from solution import count_components


def test_classic():
    assert count_components(5, [[0, 1], [1, 2], [3, 4]]) == 2


def test_all_connected():
    assert count_components(4, [[0, 1], [1, 2], [2, 3]]) == 1


def test_no_edges():
    assert count_components(4, []) == 4


def test_single_node():
    assert count_components(1, []) == 1


def test_multiple_pairs_and_isolated_nodes():
    assert count_components(6, [[0, 1], [2, 3]]) == 4
