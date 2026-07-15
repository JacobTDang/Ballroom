from solution import length_of_longest_substring


def test_length_of_longest_substring_case_01():
    assert length_of_longest_substring("abcabcbb") == 3


def test_length_of_longest_substring_case_02():
    assert length_of_longest_substring("bbbbb") == 1


def test_length_of_longest_substring_case_03():
    assert length_of_longest_substring("pwwkew") == 3


def test_length_of_longest_substring_case_04():
    assert length_of_longest_substring("") == 0


def test_length_of_longest_substring_case_05():
    assert length_of_longest_substring(" ") == 1


def test_length_of_longest_substring_case_06():
    assert length_of_longest_substring("au") == 2


def test_length_of_longest_substring_case_07():
    assert length_of_longest_substring("dvdf") == 3


def test_length_of_longest_substring_case_08():
    assert length_of_longest_substring("abba") == 2


def test_length_of_longest_substring_case_09():
    assert length_of_longest_substring("tmmzuxt") == 5


def test_length_of_longest_substring_case_10():
    assert length_of_longest_substring("aaaaaaaaaa") == 1


def test_length_of_longest_substring_case_11():
    assert length_of_longest_substring("abcdefg") == 7
