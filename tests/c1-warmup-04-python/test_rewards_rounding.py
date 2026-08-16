from solution import round_up_total


def test_example_mixed():
    assert round_up_total([250, 300, 199]) == 51


def test_example_smallest_price():
    assert round_up_total([1]) == 99


def test_example_empty():
    assert round_up_total([]) == 0


def test_exact_multiple_contributes_zero():
    assert round_up_total([100]) == 0


def test_all_exact_multiples():
    assert round_up_total([100, 200, 3000]) == 0


def test_one_below_multiple():
    assert round_up_total([99]) == 1


def test_just_over_multiple():
    assert round_up_total([101]) == 99


def test_large_price():
    assert round_up_total([99_950]) == 50
