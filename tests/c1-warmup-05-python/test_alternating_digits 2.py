from solution import is_alternating


def test_example_alternating():
    assert is_alternating("2743") is True


def test_example_ends_same_parity():
    assert is_alternating("2744") is False


def test_example_single_digit():
    assert is_alternating("7") is True


def test_pair_alternating():
    assert is_alternating("12") is True


def test_pair_both_odd():
    assert is_alternating("13") is False


def test_pair_both_even():
    assert is_alternating("08") is False


def test_break_in_middle():
    assert is_alternating("1243") is False


def test_long_alternating_even_start():
    assert is_alternating("0123456789") is True


def test_zero_alone():
    assert is_alternating("0") is True
