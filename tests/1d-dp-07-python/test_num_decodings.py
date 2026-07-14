from solution import num_decodings


def test_two_ways():
    assert num_decodings("12") == 2


def test_three_ways():
    assert num_decodings("226") == 3


def test_leading_zero():
    assert num_decodings("06") == 0


def test_two_digit_only():
    assert num_decodings("10") == 1


def test_single_digit():
    assert num_decodings("5") == 1


def test_lone_zero():
    assert num_decodings("0") == 0


def test_just_over_twenty_six():
    assert num_decodings("27") == 1


def test_unresolvable_zero_pair():
    assert num_decodings("100") == 0


def test_longer_multiple_ways():
    assert num_decodings("11106") == 2
