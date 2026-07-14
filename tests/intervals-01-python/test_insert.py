from solution import insert


def test_merge_into_middle():
    assert insert([[1, 3], [6, 9]], [2, 5]) == [[1, 5], [6, 9]]


def test_merge_several():
    assert insert([[1, 2], [3, 5], [6, 7], [8, 10], [12, 16]], [4, 8]) == [
        [1, 2],
        [3, 10],
        [12, 16],
    ]


def test_empty_list():
    assert insert([], [5, 7]) == [[5, 7]]


def test_insert_after_all_no_overlap():
    assert insert([[1, 5]], [6, 8]) == [[1, 5], [6, 8]]


def test_insert_before_all_no_overlap():
    assert insert([[1, 5]], [0, 0]) == [[0, 0], [1, 5]]


def test_new_interval_swallows_everything():
    assert insert([[2, 3], [4, 5], [6, 7]], [1, 10]) == [[1, 10]]


def test_touching_intervals_merge():
    assert insert([[1, 2]], [2, 3]) == [[1, 3]]


def test_gap_of_one_does_not_merge():
    assert insert([[1, 2]], [3, 4]) == [[1, 2], [3, 4]]


def test_identical_interval():
    assert insert([[3, 5]], [3, 5]) == [[3, 5]]


def test_new_interval_contained_within_existing():
    assert insert([[1, 10]], [3, 5]) == [[1, 10]]


def test_larger_list_merges_several_in_the_middle():
    intervals = [[1, 2], [3, 4], [5, 6], [7, 8], [9, 10], [11, 12], [13, 14], [15, 16], [17, 18], [19, 20]]
    assert insert(intervals, [6, 15]) == [[1, 2], [3, 4], [5, 16], [17, 18], [19, 20]]


def test_constraint_boundary_values():
    assert insert([[50000, 60000]], [0, 100000]) == [[0, 100000]]
