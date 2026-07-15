from solution import max_area_of_island


def test_max_area_of_island_case_1():
    grid = [
        [0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0],
        [0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0],
        [0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0],
        [0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0],
    ]
    assert max_area_of_island(grid) == 6


def test_max_area_of_island_case_2():
    assert max_area_of_island([[0, 0, 0, 0, 0, 0, 0, 0]]) == 0


def test_max_area_of_island_case_3():
    assert max_area_of_island([[1]]) == 1


def test_max_area_of_island_case_4():
    assert max_area_of_island([[1, 1], [1, 1]]) == 4


def test_max_area_of_island_case_5():
    assert max_area_of_island([[1, 0, 1], [0, 0, 0], [1, 0, 1]]) == 1
