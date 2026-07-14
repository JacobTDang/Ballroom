from solution import is_interleave


def test_classic():
    assert is_interleave("aabcc", "dbbca", "aadbbcbcac") is True


def test_not_interleaved():
    assert is_interleave("aabcc", "dbbca", "aadbbbaccc") is False


def test_all_empty():
    assert is_interleave("", "", "") is True


def test_one_empty():
    assert is_interleave("a", "", "a") is True


def test_length_mismatch():
    assert is_interleave("abc", "def", "abcde") is False


def test_first_empty_match():
    assert is_interleave("", "abc", "abc") is True


def test_first_empty_mismatch():
    assert is_interleave("", "abc", "abd") is False


def test_ambiguous_multiple_ways():
    assert is_interleave("ab", "ab", "abab") is True


def test_requires_backtrack_choice():
    assert is_interleave("ab", "ab", "aabb") is True
