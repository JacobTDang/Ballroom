from solution import longest_common_subsequence


def test_classic():
    assert longest_common_subsequence("abcde", "ace") == 3


def test_identical():
    assert longest_common_subsequence("abc", "abc") == 3


def test_no_common():
    assert longest_common_subsequence("abc", "def") == 0


def test_empty_first():
    assert longest_common_subsequence("", "abc") == 0


def test_different_order():
    assert longest_common_subsequence("abc", "acb") == 2


def test_interspersed_noise():
    assert longest_common_subsequence("aggtab", "gxtxayb") == 4


def test_repeated_chars():
    assert longest_common_subsequence("aaaa", "aa") == 2


def test_single_char_match():
    assert longest_common_subsequence("a", "a") == 1


def test_both_empty():
    assert longest_common_subsequence("", "") == 0
