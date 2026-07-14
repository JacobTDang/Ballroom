from solution import spiral_order


def test_spiral_order_3x3():
    matrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9],
    ]
    assert spiral_order(matrix) == [1, 2, 3, 6, 9, 8, 7, 4, 5]


def test_spiral_order_3x4():
    matrix = [
        [1, 2, 3, 4],
        [5, 6, 7, 8],
        [9, 10, 11, 12],
    ]
    assert spiral_order(matrix) == [1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7]


def test_spiral_order_single_row():
    assert spiral_order([[1, 2, 3, 4]]) == [1, 2, 3, 4]


def test_spiral_order_single_column():
    assert spiral_order([[1], [2], [3]]) == [1, 2, 3]


def test_spiral_order_single_element():
    assert spiral_order([[7]]) == [7]


def test_spiral_order_4x3():
    matrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9],
        [10, 11, 12],
    ]
    assert spiral_order(matrix) == [1, 2, 3, 6, 9, 12, 11, 10, 7, 4, 5, 8]


def test_spiral_order_2x2():
    matrix = [
        [1, 2],
        [3, 4],
    ]
    assert spiral_order(matrix) == [1, 2, 4, 3]


def test_spiral_order_negative_values():
    matrix = [
        [-1, -2],
        [-3, -4],
    ]
    assert spiral_order(matrix) == [-1, -2, -4, -3]


def test_spiral_order_4x4():
    matrix = [
        [1, 2, 3, 4],
        [5, 6, 7, 8],
        [9, 10, 11, 12],
        [13, 14, 15, 16],
    ]
    assert spiral_order(matrix) == [1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10]
