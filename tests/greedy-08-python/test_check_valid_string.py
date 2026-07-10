from solution import check_valid_string


def test_simple():
    assert check_valid_string("()") is True


def test_star_balances():
    assert check_valid_string("(*))") is True


def test_unbalanced():
    assert check_valid_string("(()") is False


def test_all_stars():
    assert check_valid_string("***") is True


def test_single_close():
    assert check_valid_string(")") is False
