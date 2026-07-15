from solution import check_inclusion


def test_check_inclusion_case_1():
    assert check_inclusion("ab", "eidbaooo") is True


def test_check_inclusion_case_2():
    assert check_inclusion("ab", "eidboaoo") is False


def test_check_inclusion_case_3():
    assert check_inclusion("adc", "dcda") is True


def test_check_inclusion_case_4():
    assert check_inclusion("hello", "ooolleoooleh") is False


def test_check_inclusion_case_5():
    assert check_inclusion("a", "a") is True


def test_check_inclusion_case_6():
    assert check_inclusion("abc", "ab") is False


def test_check_inclusion_case_7():
    assert check_inclusion("aa", "ab") is False


def test_check_inclusion_case_8():
    assert check_inclusion("abcd", "dcba") is True


def test_check_inclusion_case_9():
    assert check_inclusion("ab", "a") is False
