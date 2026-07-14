from solution import unique_paths


def test_classic():
    assert unique_paths(3, 7) == 28


def test_small():
    assert unique_paths(3, 2) == 3


def test_single_cell():
    assert unique_paths(1, 1) == 1


def test_single_row():
    assert unique_paths(1, 5) == 1


def test_single_column():
    assert unique_paths(5, 1) == 1


def test_square():
    assert unique_paths(3, 3) == 6


def test_larger_square():
    assert unique_paths(10, 10) == 48620


def test_boundary_max_with_min_other():
    assert unique_paths(2, 100) == 100
