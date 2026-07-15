from solution import find_median_sorted_arrays


def test_find_median_sorted_arrays_case_1():
    assert find_median_sorted_arrays([1, 3], [2]) == 2.0


def test_find_median_sorted_arrays_case_2():
    assert find_median_sorted_arrays([1, 2], [3, 4]) == 2.5


def test_find_median_sorted_arrays_case_3():
    assert find_median_sorted_arrays([], [1]) == 1.0


def test_find_median_sorted_arrays_case_4():
    assert find_median_sorted_arrays([2], []) == 2.0


def test_find_median_sorted_arrays_case_5():
    assert find_median_sorted_arrays([1, 2, 3, 4, 5], [6, 7, 8, 9, 10]) == 5.5


def test_find_median_sorted_arrays_case_6():
    assert find_median_sorted_arrays([1, 2, 3], [4, 5, 6, 7, 8, 9]) == 5.0


def test_find_median_sorted_arrays_case_7():
    assert find_median_sorted_arrays([1], [2, 3, 4, 5]) == 3.0


def test_find_median_sorted_arrays_case_8():
    assert find_median_sorted_arrays([100, 200], [1, 2, 3]) == 3.0
