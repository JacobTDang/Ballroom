from solution import change


def test_classic():
    assert change(5, [1, 2, 5]) == 4


def test_no_way():
    assert change(3, [2]) == 0


def test_zero_amount():
    assert change(0, [1, 2, 3]) == 1


def test_exact_single_coin():
    assert change(10, [10]) == 1


def test_larger_amount():
    assert change(10, [1, 2, 5]) == 10


def test_single_coin_no_divide():
    assert change(7, [3]) == 0


def test_more_denominations():
    assert change(10, [2, 5, 3, 6]) == 5


def test_boundary_amount_single_coin():
    assert change(500, [1]) == 1
