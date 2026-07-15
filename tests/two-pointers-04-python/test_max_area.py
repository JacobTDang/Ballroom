from solution import max_area


def test_max_area_case_01():
    assert max_area([1, 8, 6, 2, 5, 4, 8, 3, 7]) == 49


def test_max_area_case_02():
    assert max_area([1, 1]) == 1


def test_max_area_case_03():
    assert max_area([4, 3, 2, 1, 4]) == 16


def test_max_area_case_04():
    assert max_area([1, 2, 1]) == 2


def test_max_area_case_05():
    assert max_area([1, 2, 4, 3]) == 4


def test_max_area_case_06():
    assert max_area([1, 3, 2, 5, 25, 24, 5]) == 24


def test_max_area_case_07():
    assert max_area([0, 2]) == 0


def test_max_area_case_08():
    assert max_area([5, 5, 5, 5, 5]) == 20


def test_max_area_case_09():
    assert max_area([1, 2, 3, 4, 5, 25, 1]) == 9


def test_max_area_case_10():
    assert max_area([10000, 10000]) == 10000
