from solution import merge


def test_overlapping_pair():
    assert merge([[1, 3], [2, 6], [8, 10], [15, 18]]) == [[1, 6], [8, 10], [15, 18]]


def test_touching_endpoints_merge():
    assert merge([[1, 4], [4, 5]]) == [[1, 5]]


def test_unsorted_input():
    assert merge([[15, 18], [2, 6], [1, 3], [8, 10]]) == [[1, 6], [8, 10], [15, 18]]


def test_single_interval():
    assert merge([[1, 4]]) == [[1, 4]]


def test_one_interval_fully_contains_another():
    assert merge([[1, 10], [2, 3], [4, 5]]) == [[1, 10]]


def test_no_overlaps_at_all():
    assert merge([[1, 2], [3, 4], [5, 6]]) == [[1, 2], [3, 4], [5, 6]]


def test_boundary_values_no_overlap():
    assert merge([[0, 1], [9999, 10000]]) == [[0, 1], [9999, 10000]]


def test_larger_input_multiple_merge_chains():
    intervals = [[1, 3], [2, 4], [5, 7], [6, 8], [10, 12], [15, 20], [18, 25], [30, 31]]
    assert merge(intervals) == [[1, 4], [5, 8], [10, 12], [15, 25], [30, 31]]
