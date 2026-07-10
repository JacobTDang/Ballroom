from solution import can_complete_circuit


def test_classic():
    assert can_complete_circuit([1, 2, 3, 4, 5], [3, 4, 5, 1, 2]) == 3


def test_impossible():
    assert can_complete_circuit([2, 3, 4], [3, 4, 3]) == -1


def test_single_exact():
    assert can_complete_circuit([5], [4]) == 0


def test_single_insufficient():
    assert can_complete_circuit([3], [4]) == -1


def test_start_at_last_index():
    assert can_complete_circuit([5, 1, 2, 3, 4], [4, 4, 1, 5, 1]) == 4
