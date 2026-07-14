from solution import rotate_image


def test_rotate_3x3():
    matrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9],
    ]
    rotate_image(matrix)
    assert matrix == [
        [7, 4, 1],
        [8, 5, 2],
        [9, 6, 3],
    ]


def test_rotate_2x2():
    matrix = [
        [1, 2],
        [3, 4],
    ]
    rotate_image(matrix)
    assert matrix == [
        [3, 1],
        [4, 2],
    ]


def test_rotate_1x1():
    matrix = [[5]]
    rotate_image(matrix)
    assert matrix == [[5]]


def test_rotate_4x4():
    matrix = [
        [1, 2, 3, 4],
        [5, 6, 7, 8],
        [9, 10, 11, 12],
        [13, 14, 15, 16],
    ]
    rotate_image(matrix)
    assert matrix == [
        [13, 9, 5, 1],
        [14, 10, 6, 2],
        [15, 11, 7, 3],
        [16, 12, 8, 4],
    ]


def test_negative_values():
    matrix = [
        [-1, -2],
        [-3, -4],
    ]
    rotate_image(matrix)
    assert matrix == [
        [-3, -1],
        [-4, -2],
    ]


def test_all_same_values():
    matrix = [
        [7, 7, 7],
        [7, 7, 7],
        [7, 7, 7],
    ]
    rotate_image(matrix)
    assert matrix == [
        [7, 7, 7],
        [7, 7, 7],
        [7, 7, 7],
    ]


def test_with_zero():
    matrix = [
        [0, 1],
        [2, 3],
    ]
    rotate_image(matrix)
    assert matrix == [
        [2, 0],
        [3, 1],
    ]


def test_5x5():
    matrix = [
        [1, 2, 3, 4, 5],
        [6, 7, 8, 9, 10],
        [11, 12, 13, 14, 15],
        [16, 17, 18, 19, 20],
        [21, 22, 23, 24, 25],
    ]
    rotate_image(matrix)
    assert matrix == [
        [21, 16, 11, 6, 1],
        [22, 17, 12, 7, 2],
        [23, 18, 13, 8, 3],
        [24, 19, 14, 9, 4],
        [25, 20, 15, 10, 5],
    ]
