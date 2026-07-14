from solution import word_break


def test_classic():
    assert word_break("leetcode", ["leet", "code"]) is True


def test_reused_word():
    assert word_break("applepenapple", ["apple", "pen"]) is True


def test_impossible():
    assert word_break("catsandog", ["cats", "dog", "sand", "and", "cat"]) is False


def test_single_word():
    assert word_break("a", ["a"]) is True


def test_leftover_char_unmatched():
    assert word_break("ab", ["a"]) is False


def test_trailing_char_never_matches():
    assert word_break("aaaaaaaaaaaaaaaaaaaab", ["a", "aa"]) is False


def test_multiple_paths():
    assert word_break("cars", ["car", "ca", "rs"]) is True


def test_simple_concatenation():
    assert word_break("goalspecial", ["goal", "special"]) is True
