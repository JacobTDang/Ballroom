from solution import is_match


def test_no_star_mismatch():
    assert is_match("aa", "a") is False


def test_star_repeat():
    assert is_match("aa", "a*") is True


def test_classic():
    assert is_match("aab", "c*a*b") is True


def test_long_no_match():
    assert is_match("mississippi", "mis*is*p*.") is False


def test_both_empty():
    assert is_match("", "") is True


def test_empty_string_star_zero():
    assert is_match("", "a*") is True


def test_dot_matches_any():
    assert is_match("a", ".") is True


def test_dot_star_matches_all():
    assert is_match("ab", ".*") is True


def test_longer_match():
    assert is_match("mississippi", "mis*is*ip*.") is True


def test_dot_star_trailing_literal_fails():
    assert is_match("ab", ".*c") is False
