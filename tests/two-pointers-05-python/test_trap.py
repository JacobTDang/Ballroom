from solution import trap


def test_trap_case_01():
    assert trap([0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]) == 6


def test_trap_case_02():
    assert trap([4, 2, 0, 3, 2, 5]) == 9


def test_trap_case_03():
    assert trap([]) == 0


def test_trap_case_04():
    assert trap([1, 2, 3, 4, 5]) == 0


def test_trap_case_05():
    assert trap([5, 4, 3, 2, 1]) == 0


def test_trap_case_06():
    assert trap([3, 0, 3]) == 3


def test_trap_case_07():
    assert trap([2, 0, 2]) == 2


def test_trap_case_08():
    assert trap([5]) == 0


def test_trap_case_09():
    assert trap([1, 0, 1]) == 1


def test_trap_case_10():
    assert trap([4, 4, 4, 4]) == 0


def test_trap_case_11():
    assert trap([5, 2, 1, 2, 1, 5]) == 14
