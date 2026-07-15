from solution import oranges_rotting


def test_oranges_rotting_case_1():
    assert oranges_rotting([[2, 1, 1], [1, 1, 0], [0, 1, 1]]) == 4


def test_oranges_rotting_case_2():
    assert oranges_rotting([[2, 1, 1], [0, 1, 1], [1, 0, 1]]) == -1


def test_oranges_rotting_case_3():
    assert oranges_rotting([[0, 2]]) == 0


def test_oranges_rotting_case_4():
    assert oranges_rotting([[0]]) == 0


def test_oranges_rotting_case_5():
    assert oranges_rotting([[2, 1, 1], [1, 1, 1], [1, 1, 2]]) == 2


def test_oranges_rotting_case_6():
    assert oranges_rotting([[2, 2], [2, 2]]) == 0
