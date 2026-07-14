from solution import longest_palindrome


def test_odd_tie():
    assert longest_palindrome("babad") == "bab"


def test_even():
    assert longest_palindrome("cbbd") == "bb"


def test_single_char():
    assert longest_palindrome("a") == "a"


def test_whole_string():
    assert longest_palindrome("abba") == "abba"


def test_all_same_longer():
    assert longest_palindrome("aaaaa") == "aaaaa"


def test_no_repeat():
    assert longest_palindrome("abcde") == "a"


def test_buried_in_larger_string():
    assert longest_palindrome("zzabcbayy") == "abcba"
