from solution import rob_circular


def test_classic():
    assert rob_circular([2, 3, 2]) == 3


def test_four_houses():
    assert rob_circular([1, 2, 3, 1]) == 4


def test_three_in_a_row():
    assert rob_circular([1, 2, 3]) == 3


def test_single_house():
    assert rob_circular([5]) == 5


def test_two_houses():
    assert rob_circular([5, 10]) == 10


def test_all_zeros():
    assert rob_circular([0, 0, 0, 0]) == 0


def test_larger_alternating():
    assert rob_circular([2, 3, 2, 3, 2, 3, 2]) == 9


def test_boundary_max_values():
    assert rob_circular([1000, 1000, 1000]) == 1000
