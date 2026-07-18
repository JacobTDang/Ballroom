from solution import count_paths


def test_count_paths_two_by_two_open():
    assert count_paths([[0, 0], [0, 0]]) == 2


def test_count_paths_trivial_single_cell():
    assert count_paths([[0]]) == 1


def test_count_paths_destination_blocked():
    assert count_paths([[0, 0], [0, 1]]) == 0


def test_count_paths_fully_blocked():
    # Anti-overfit: the destination cell itself is open, but every
    # route to it is blocked. A fix that just hardcodes the
    # destination base case to 1 without preserving the blocked-cell
    # check must still get 0 here.
    assert count_paths([[0, 1], [1, 0]]) == 0


def test_count_paths_with_obstacle():
    assert count_paths([[0, 0, 0], [0, 1, 0], [0, 0, 0]]) == 2


def test_count_paths_single_row():
    assert count_paths([[0, 0, 0]]) == 1
