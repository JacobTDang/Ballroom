from solution import length_of_longest_substring


def test_length_of_longest_substring():
    assert length_of_longest_substring("abcabcbb") == 3
    assert length_of_longest_substring("bbbbb") == 1
    assert length_of_longest_substring("pwwkew") == 3
    assert length_of_longest_substring("") == 0
    assert length_of_longest_substring(" ") == 1
    assert length_of_longest_substring("au") == 2
    assert length_of_longest_substring("dvdf") == 3
    assert length_of_longest_substring("abba") == 2
    assert length_of_longest_substring("tmmzuxt") == 5
    assert length_of_longest_substring("aaaaaaaaaa") == 1
    assert length_of_longest_substring("abcdefg") == 7
