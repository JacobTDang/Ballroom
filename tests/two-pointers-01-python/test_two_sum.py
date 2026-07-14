from solution import two_sum


def test_two_sum():
    assert two_sum([2, 7, 11, 15], 9) == [1, 2]
    assert two_sum([2, 3, 4], 6) == [1, 3]
    assert two_sum([-1, 0], -1) == [1, 2]
    assert two_sum([3, 3], 6) == [1, 2]
    assert two_sum([1, 2, 3, 4, 4, 9, 56, 90], 8) == [4, 5]
    assert two_sum([-8, -3, 0, 4, 9, 13], 5) == [1, 6]
    assert two_sum([-8, -3, 0, 4, 9, 13], 4) == [3, 4]
    assert two_sum([-3, -1, 0, 2, 4, 5], 1) == [1, 5]
    assert two_sum([5, 25, 75], 100) == [2, 3]
