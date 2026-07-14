from solution import largest_rectangle_area


def test_largest_rectangle_area():
    assert largest_rectangle_area([2, 1, 5, 6, 2, 3]) == 10
    assert largest_rectangle_area([2, 4]) == 4
    assert largest_rectangle_area([1]) == 1
    assert largest_rectangle_area([0, 0]) == 0
    assert largest_rectangle_area([5, 5, 5, 5]) == 20
    assert largest_rectangle_area([5, 4, 3, 2, 1]) == 9
    assert largest_rectangle_area([1, 2, 3, 4, 5]) == 9
