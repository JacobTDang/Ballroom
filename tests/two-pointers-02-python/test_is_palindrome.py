from solution import is_palindrome


def test_is_palindrome_case_01():
    assert is_palindrome("A man, a plan, a canal: Panama") is True


def test_is_palindrome_case_02():
    assert is_palindrome("race a car") is False


def test_is_palindrome_case_03():
    assert is_palindrome(" ") is True


def test_is_palindrome_case_04():
    assert is_palindrome("0P") is False


def test_is_palindrome_case_05():
    assert is_palindrome("Was it a car or a cat I saw?") is True


def test_is_palindrome_case_06():
    assert is_palindrome(".,") is True


def test_is_palindrome_case_07():
    assert is_palindrome("a_b") is False


def test_is_palindrome_case_08():
    assert is_palindrome("12321") is True


def test_is_palindrome_case_09():
    assert is_palindrome("ab") is False


def test_is_palindrome_case_10():
    assert is_palindrome("") is True


def test_is_palindrome_case_11():
    assert is_palindrome("Able was I, ere I saw Elba") is True
