from solution import can_finish


def test_no_cycle():
    assert can_finish(2, [[1, 0]]) is True


def test_cycle():
    assert can_finish(2, [[1, 0], [0, 1]]) is False


def test_no_prerequisites():
    assert can_finish(5, []) is True


def test_longer_cycle():
    assert can_finish(4, [[1, 0], [2, 1], [3, 2], [0, 3]]) is False
