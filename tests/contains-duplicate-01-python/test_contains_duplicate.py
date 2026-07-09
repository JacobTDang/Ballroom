from solution import contains_duplicate


def test_contains_duplicate():
    assert contains_duplicate([1, 2, 3, 1]) is True
    assert contains_duplicate([1, 2, 3, 4]) is False
    assert contains_duplicate([1, 1, 1, 3, 3, 4, 3, 2, 4, 2]) is True
    assert contains_duplicate([1]) is False
    assert contains_duplicate([-1, -1]) is True
    assert contains_duplicate([0, 4, 5, 0, 3, 6]) is True
