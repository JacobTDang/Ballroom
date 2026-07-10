from solution import is_anagram


def test_is_anagram():
    assert is_anagram("anagram", "nagaram") is True
    assert is_anagram("rat", "car") is False
    assert is_anagram("ab", "a") is False
    assert is_anagram("aacc", "ccac") is False
    assert is_anagram("a", "a") is True
    assert is_anagram("aabbcc", "abcabc") is True
