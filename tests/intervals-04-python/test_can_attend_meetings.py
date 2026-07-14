from solution import can_attend_meetings


def test_overlap_in_the_middle():
    assert can_attend_meetings([[0, 30], [5, 10], [15, 20]]) is False


def test_no_overlap():
    assert can_attend_meetings([[7, 10], [2, 4]]) is True


def test_touching_endpoints_is_fine():
    assert can_attend_meetings([[5, 10], [10, 15]]) is True


def test_empty_schedule():
    assert can_attend_meetings([]) is True


def test_single_meeting():
    assert can_attend_meetings([[3, 8]]) is True


def test_unsorted_overlap_at_the_end():
    assert can_attend_meetings([[13, 15], [1, 5], [6, 8], [14, 20]]) is False


def test_boundary_values_overlap_by_one():
    assert can_attend_meetings([[0, 1000000], [999999, 1000000]]) is False


def test_larger_schedule_all_sequential():
    intervals = [[0, 10], [10, 20], [20, 30], [30, 40], [40, 50], [50, 60], [60, 70]]
    assert can_attend_meetings(intervals) is True
