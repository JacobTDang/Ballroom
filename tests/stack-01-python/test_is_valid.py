from solution import is_valid


def test_is_valid():
    assert is_valid("()") is True
    assert is_valid("()[]{}") is True
    assert is_valid("(]") is False
    assert is_valid("([)]") is False
    assert is_valid("{[]}") is True
    assert is_valid("(") is False
    assert is_valid("]") is False
    assert is_valid("") is True
    assert is_valid("((()))") is True
    assert is_valid("(((") is False
    assert is_valid("){") is False
