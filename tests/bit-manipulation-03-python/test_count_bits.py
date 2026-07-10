from solution import count_bits


def test_classic():
    assert count_bits(4) == [0, 1, 1, 2, 1]


def test_small():
    assert count_bits(2) == [0, 1, 1]


def test_zero():
    assert count_bits(0) == [0]


def test_larger():
    assert count_bits(15) == [0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4]
