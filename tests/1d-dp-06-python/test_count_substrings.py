from solution import count_substrings


def test_three_distinct():
    assert count_substrings("abc") == 3


def test_all_same():
    assert count_substrings("aaa") == 6


def test_odd_palindrome():
    assert count_substrings("aba") == 4


def test_single_char():
    assert count_substrings("z") == 1


def test_two_same():
    assert count_substrings("aa") == 3


def test_two_different():
    assert count_substrings("ab") == 2


def test_larger_all_same():
    assert count_substrings("aaaaa") == 15


def test_nested_palindromes():
    assert count_substrings("aabaa") == 9
