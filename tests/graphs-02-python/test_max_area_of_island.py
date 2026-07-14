from solution import max_area_of_island


def test_max_area_of_island():
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
    assert max_area_of_island([[0, 0, 0, 0, 0, 0, 0, 0]]) == 0
    assert max_area_of_island([[1]]) == 1
    assert max_area_of_island([[1, 1], [1, 1]]) == 4
    assert max_area_of_island([[1, 0, 1], [0, 0, 0], [1, 0, 1]]) == 1
