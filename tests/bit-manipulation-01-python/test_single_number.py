from solution import single_number


def test_classic():
    assert single_number([2, 2, 1]) == 1


def test_longer_mix():
    assert single_number([4, 1, 2, 1, 2]) == 4


def test_single_element():
    assert single_number([7]) == 7


def test_negative_numbers():
    assert single_number([-1, -1, -2]) == -2


def test_boundary_values():
    assert single_number([1000, 1000, -1000]) == -1000


def test_larger_mixed_set():
    assert single_number([5, 3, 5, 4, 3]) == 4
