from solution import erase_overlap_intervals


def test_one_overlap_to_remove():
    assert erase_overlap_intervals([[1, 2], [2, 3], [3, 4], [1, 3]]) == 1


def test_all_identical():
    assert erase_overlap_intervals([[1, 2], [1, 2], [1, 2]]) == 2


def test_touching_endpoints_no_removal():
    assert erase_overlap_intervals([[1, 2], [2, 3]]) == 0


def test_single_interval():
    assert erase_overlap_intervals([[1, 2]]) == 0


def test_already_non_overlapping():
    assert erase_overlap_intervals([[1, 2], [3, 4], [5, 6]]) == 0


def test_heavy_overlap_needs_two_removals():
    assert erase_overlap_intervals([[1, 100], [11, 22], [1, 11], [2, 12]]) == 2


def test_boundary_values_touching_not_overlapping():
    assert erase_overlap_intervals([[-50000, -49999], [-49999, 50000]]) == 0


def test_all_sharing_start_most_must_go():
    assert erase_overlap_intervals([[1, 2], [1, 3], [1, 4], [1, 5], [1, 6]]) == 4
