from solution import reverse_bits


def test_one():
    assert reverse_bits(1) == 2147483648


def test_zero():
    assert reverse_bits(0) == 0


def test_all_ones():
    assert reverse_bits(4294967295) == 4294967295


def test_two():
    assert reverse_bits(2) == 1073741824


def test_classic():
    assert reverse_bits(43261596) == 964176192
