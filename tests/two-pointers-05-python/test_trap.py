from solution import trap


def test_trap():
    assert trap([0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]) == 6
    assert trap([4, 2, 0, 3, 2, 5]) == 9
    assert trap([]) == 0
    assert trap([1, 2, 3, 4, 5]) == 0
    assert trap([5, 4, 3, 2, 1]) == 0
    assert trap([3, 0, 3]) == 3
    assert trap([2, 0, 2]) == 2
    assert trap([5]) == 0
    assert trap([1, 0, 1]) == 1
    assert trap([4, 4, 4, 4]) == 0
    assert trap([5, 2, 1, 2, 1, 5]) == 14
