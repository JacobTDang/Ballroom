from solution import largest_rectangle_area


def test_largest_rectangle_area_case_1():
    assert largest_rectangle_area([2, 1, 5, 6, 2, 3]) == 10


def test_largest_rectangle_area_case_2():
    assert largest_rectangle_area([2, 4]) == 4


def test_largest_rectangle_area_case_3():
    assert largest_rectangle_area([1]) == 1


def test_largest_rectangle_area_case_4():
    assert largest_rectangle_area([0, 0]) == 0


def test_largest_rectangle_area_case_5():
    assert largest_rectangle_area([5, 5, 5, 5]) == 20


def test_largest_rectangle_area_case_6():
    assert largest_rectangle_area([5, 4, 3, 2, 1]) == 9


def test_largest_rectangle_area_case_7():
    assert largest_rectangle_area([1, 2, 3, 4, 5]) == 9
