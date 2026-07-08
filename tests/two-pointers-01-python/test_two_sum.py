from solution import two_sum


def test_two_sum():
    assert two_sum([2, 7, 11, 15], 9) == [1, 2]
    assert two_sum([2, 3, 4], 6) == [1, 3]
    assert two_sum([-1, 0], -1) == [1, 2]
