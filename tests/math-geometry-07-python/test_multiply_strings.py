from solution import multiply_strings


def test_multiply_strings_single_digits():
    assert multiply_strings("2", "3") == "6"


def test_multiply_strings_multi_digits():
    assert multiply_strings("123", "456") == "56088"


def test_multiply_strings_zero_first():
    assert multiply_strings("0", "12345") == "0"


def test_multiply_strings_all_nines():
    assert multiply_strings("999", "999") == "998001"


def test_multiply_strings_one_by_one():
    assert multiply_strings("1", "1") == "1"


def test_multiply_strings_both_zero():
    assert multiply_strings("0", "0") == "0"


def test_multiply_strings_identity_by_larger():
    assert multiply_strings("1", "999") == "999"


def test_multiply_strings_larger_multi_digit():
    assert multiply_strings("12345", "67890") == "838102050"


def test_multiply_strings_single_digit_by_multi_digit():
    assert multiply_strings("9", "123") == "1107"
