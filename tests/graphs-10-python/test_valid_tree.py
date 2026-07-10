from solution import valid_tree


def test_valid():
    assert valid_tree(5, [[0, 1], [0, 2], [0, 3], [1, 4]]) is True


def test_has_cycle():
    assert valid_tree(5, [[0, 1], [1, 2], [2, 3], [1, 3], [1, 4]]) is False


def test_disconnected():
    assert valid_tree(4, [[0, 1], [2, 3]]) is False


def test_single_node():
    assert valid_tree(1, []) is True
