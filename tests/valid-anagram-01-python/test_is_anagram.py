from solution import is_anagram


def test_is_anagram_case_01():
    assert is_anagram("anagram", "nagaram") is True


def test_is_anagram_case_02():
    assert is_anagram("rat", "car") is False


def test_is_anagram_case_03():
    assert is_anagram("ab", "a") is False


def test_is_anagram_case_04():
    assert is_anagram("aacc", "ccac") is False


def test_is_anagram_case_05():
    assert is_anagram("a", "a") is True


def test_is_anagram_case_06():
    assert is_anagram("aabbcc", "abcabc") is True


def test_is_anagram_case_07():
    assert is_anagram("listen", "silent") is True


def test_is_anagram_case_08():
    assert is_anagram("aaab", "aabb") is False


def test_is_anagram_case_09():
    assert is_anagram("a", "b") is False


def test_is_anagram_case_10():
    assert is_anagram("abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcba") is True


def test_is_anagram_case_11():
    assert is_anagram("aaaa", "aaaa") is True
