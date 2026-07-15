from solution import contains_duplicate


def test_contains_duplicate_case_01():
    assert contains_duplicate([1, 2, 3, 1]) is True


def test_contains_duplicate_case_02():
    assert contains_duplicate([1, 2, 3, 4]) is False


def test_contains_duplicate_case_03():
    assert contains_duplicate([1, 1, 1, 3, 3, 4, 3, 2, 4, 2]) is True


def test_contains_duplicate_case_04():
    assert contains_duplicate([1]) is False


def test_contains_duplicate_case_05():
    assert contains_duplicate([1, 1]) is True


def test_contains_duplicate_case_06():
    assert contains_duplicate([1, 2]) is False


def test_contains_duplicate_case_07():
    assert contains_duplicate([-1, -1]) is True


def test_contains_duplicate_case_08():
    assert contains_duplicate([-5, -3, -1, 1, 3, 5]) is False


def test_contains_duplicate_case_09():
    assert contains_duplicate([0, 4, 5, 0, 3, 6]) is True


def test_contains_duplicate_case_10():
    assert contains_duplicate([7, 7, 7, 7, 7]) is True


def test_contains_duplicate_case_11():
    assert contains_duplicate([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1]) is True


def test_contains_duplicate_case_12():
    assert contains_duplicate([-1000000000, 1000000000]) is False


def test_contains_duplicate_case_13():
    assert contains_duplicate([1000000000, 1000000000]) is True
