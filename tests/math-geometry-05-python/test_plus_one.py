from solution import plus_one


def test_plus_one_simple():
    assert plus_one([1, 2, 3]) == [1, 2, 4]


def test_plus_one_all_nines():
    assert plus_one([9, 9, 9]) == [1, 0, 0, 0]


def test_plus_one_single_zero():
    assert plus_one([0]) == [1]


def test_plus_one_trailing_nine():
    assert plus_one([1, 2, 9]) == [1, 3, 0]


def test_plus_one_single_nine():
    assert plus_one([9]) == [1, 0]


def test_plus_one_partial_trailing_nines():
    assert plus_one([1, 9, 9]) == [2, 0, 0]


def test_plus_one_mixed_no_carry_past_stop():
    assert plus_one([9, 8, 9, 9]) == [9, 9, 0, 0]


def test_plus_one_single_digit_not_nine():
    assert plus_one([5]) == [6]


def test_plus_one_larger_all_nines():
    assert plus_one([9, 9, 9, 9, 9]) == [1, 0, 0, 0, 0, 0]
