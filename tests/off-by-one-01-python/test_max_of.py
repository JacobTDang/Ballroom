from solution import max_of


def test_max_of():
    assert max_of([3, 1, 4, 1, 5, 9, 2, 6]) == 9
    assert max_of([-5, -1, -10]) == -1
    assert max_of([42]) == 42
    assert max_of([5, 5, 5]) == 5
    assert max_of([1, 2, 3, 4, 5, 100]) == 100
    assert max_of([-1, -2, -3, -100]) == -1
