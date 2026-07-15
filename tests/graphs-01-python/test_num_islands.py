from solution import num_islands


def grid_of(rows):
    return [list(r) for r in rows]


def test_num_islands_case_1():
    assert num_islands(grid_of(["11110", "11010", "11000", "00000"])) == 1


def test_num_islands_case_2():
    assert num_islands(grid_of(["11000", "11000", "00100", "00011"])) == 3


def test_num_islands_case_3():
    assert num_islands(grid_of(["0"])) == 0


def test_num_islands_case_4():
    assert num_islands(grid_of(["1"])) == 1


def test_num_islands_case_5():
    assert num_islands(grid_of(["000", "000"])) == 0


def test_num_islands_case_6():
    assert num_islands(grid_of(["11", "11"])) == 1
