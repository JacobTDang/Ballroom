from solution import rob


def test_classic():
    assert rob([1, 2, 3, 1]) == 4


def test_larger():
    assert rob([2, 7, 9, 3, 1]) == 12


def test_single_house():
    assert rob([5]) == 5


def test_two_houses():
    assert rob([2, 1]) == 2


def test_all_zeros():
    assert rob([0, 0, 0, 0]) == 0


def test_larger_mixed_values():
    assert rob([5, 5, 10, 100, 10, 5]) == 110


def test_boundary_max_values():
    assert rob([1000, 1000]) == 1000
