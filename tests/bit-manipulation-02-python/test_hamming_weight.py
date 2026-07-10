from solution import hamming_weight


def test_classic():
    assert hamming_weight(11) == 3


def test_zero():
    assert hamming_weight(0) == 0


def test_all_ones():
    assert hamming_weight(4294967295) == 32


def test_power_of_two():
    assert hamming_weight(1 << 31) == 1
