from solution import max_sliding_window


def test_max_sliding_window():
    assert max_sliding_window([1, 3, -1, -3, 5, 3, 6, 7], 3) == [3, 3, 5, 5, 6, 7]
    assert max_sliding_window([1], 1) == [1]
    assert max_sliding_window([1, -1], 1) == [1, -1]
    assert max_sliding_window([9, 11], 2) == [11]
    assert max_sliding_window([4, -2], 2) == [4]
    assert max_sliding_window([1, 3, 1, 2, 0, 5], 3) == [3, 3, 2, 5]
    assert max_sliding_window([7, 2, 4], 2) == [7, 4]
    assert max_sliding_window([1, 2, 3, 4, 5], 5) == [5]
    assert max_sliding_window([-7, -8, 7, 5, 7, 1, 6, 0], 4) == [7, 7, 7, 7, 7]
