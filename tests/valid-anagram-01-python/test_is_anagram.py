from solution import is_anagram


def test_is_anagram():
    assert is_anagram("anagram", "nagaram") is True
    assert is_anagram("rat", "car") is False
    assert is_anagram("ab", "a") is False
    assert is_anagram("aacc", "ccac") is False
    assert is_anagram("a", "a") is True
    assert is_anagram("aabbcc", "abcabc") is True
    assert is_anagram("listen", "silent") is True
    assert is_anagram("aaab", "aabb") is False
    assert is_anagram("a", "b") is False
    assert is_anagram("abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcba") is True
    assert is_anagram("aaaa", "aaaa") is True
