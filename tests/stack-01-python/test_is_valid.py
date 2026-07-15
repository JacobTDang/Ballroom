from solution import is_valid


def test_is_valid_case_01():
    assert is_valid("()") is True


def test_is_valid_case_02():
    assert is_valid("()[]{}") is True


def test_is_valid_case_03():
    assert is_valid("(]") is False


def test_is_valid_case_04():
    assert is_valid("([)]") is False


def test_is_valid_case_05():
    assert is_valid("{[]}") is True


def test_is_valid_case_06():
    assert is_valid("(") is False


def test_is_valid_case_07():
    assert is_valid("]") is False


def test_is_valid_case_08():
    assert is_valid("") is True


def test_is_valid_case_09():
    assert is_valid("((()))") is True


def test_is_valid_case_10():
    assert is_valid("(((") is False


def test_is_valid_case_11():
    assert is_valid("){") is False
