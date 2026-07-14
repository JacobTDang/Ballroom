from solution import min_meeting_rooms


def test_classic_two_rooms():
    assert min_meeting_rooms([[0, 30], [5, 10], [15, 20]]) == 2


def test_no_overlap_one_room():
    assert min_meeting_rooms([[7, 10], [2, 4]]) == 1


def test_touching_endpoints_share_a_room():
    assert min_meeting_rooms([[5, 10], [10, 15]]) == 1


def test_empty_schedule():
    assert min_meeting_rooms([]) == 0


def test_three_identical_meetings_need_three_rooms():
    assert min_meeting_rooms([[1, 2], [1, 2], [1, 2]]) == 3


def test_single_meeting():
    assert min_meeting_rooms([[1, 5]]) == 1


def test_five_fully_overlapping_meetings_need_five_rooms():
    assert min_meeting_rooms([[1, 100], [1, 100], [1, 100], [1, 100], [1, 100]]) == 5


def test_staggered_overlaps_larger_input():
    intervals = [[1, 10], [2, 7], [3, 19], [8, 12], [10, 20], [11, 30]]
    assert min_meeting_rooms(intervals) == 4
