from solution import is_happy


def test_is_happy_19():
    assert is_happy(19) is True


def test_is_happy_2():
    assert is_happy(2) is False


def test_is_happy_1():
    assert is_happy(1) is True


def test_is_happy_7():
    assert is_happy(7) is True


def test_is_happy_4():
    assert is_happy(4) is False


def test_is_happy_100():
    assert is_happy(100) is True


def test_is_happy_3():
    assert is_happy(3) is False


def test_is_happy_986():
    assert is_happy(986) is False


def test_is_happy_boundary_large():
    assert is_happy(2147483647) is False
