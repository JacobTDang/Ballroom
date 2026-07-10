from solution import check_inclusion


def test_check_inclusion():
    assert check_inclusion("ab", "eidbaooo") is True
    assert check_inclusion("ab", "eidboaoo") is False
    assert check_inclusion("adc", "dcda") is True
    assert check_inclusion("hello", "ooolleoooleh") is False
    assert check_inclusion("a", "a") is True
    assert check_inclusion("abc", "ab") is False
