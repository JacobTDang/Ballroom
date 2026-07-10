from solution import find_median_sorted_arrays


def test_find_median_sorted_arrays():
    assert find_median_sorted_arrays([1, 3], [2]) == 2.0
    assert find_median_sorted_arrays([1, 2], [3, 4]) == 2.5
    assert find_median_sorted_arrays([], [1]) == 1.0
    assert find_median_sorted_arrays([2], []) == 2.0
    assert find_median_sorted_arrays([1, 2, 3, 4, 5], [6, 7, 8, 9, 10]) == 5.5
