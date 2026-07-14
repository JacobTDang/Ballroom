from solution import set_zeroes


def test_set_zeroes_classic():
    matrix = [
        [1, 1, 1],
        [1, 0, 1],
        [1, 1, 1],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [1, 0, 1],
        [0, 0, 0],
        [1, 0, 1],
    ]


def test_set_zeroes_two_zeroes():
    matrix = [
        [0, 1, 2, 0],
        [3, 4, 5, 2],
        [1, 3, 1, 5],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [0, 0, 0, 0],
        [0, 4, 5, 0],
        [0, 3, 1, 0],
    ]


def test_set_zeroes_single_zero():
    matrix = [
        [1, 0],
        [1, 1],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [0, 0],
        [1, 0],
    ]


def test_set_zeroes_no_zero():
    matrix = [
        [1, 2],
        [3, 4],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [1, 2],
        [3, 4],
    ]


def test_set_zeroes_all_zeros():
    matrix = [
        [0, 0],
        [0, 0],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [0, 0],
        [0, 0],
    ]


def test_set_zeroes_corner_zero():
    matrix = [
        [0, 1],
        [1, 1],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [0, 0],
        [0, 1],
    ]


def test_set_zeroes_single_row():
    matrix = [[1, 0, 3]]
    set_zeroes(matrix)
    assert matrix == [[0, 0, 0]]


def test_set_zeroes_single_column():
    matrix = [[1], [0], [3]]
    set_zeroes(matrix)
    assert matrix == [[0], [0], [0]]


def test_set_zeroes_negative_values():
    matrix = [
        [-1, 0],
        [-2, -3],
    ]
    set_zeroes(matrix)
    assert matrix == [
        [0, 0],
        [-2, 0],
    ]
