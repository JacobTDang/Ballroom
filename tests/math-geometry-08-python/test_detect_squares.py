from solution import DetectSquares


def test_detect_squares_classic():
    ds = DetectSquares()
    ds.add([3, 10])
    ds.add([11, 2])
    ds.add([3, 2])

    assert ds.count([11, 10]) == 1
    assert ds.count([14, 8]) == 0

    ds.add([11, 2])
    assert ds.count([11, 10]) == 2


def test_detect_squares_empty_state():
    ds = DetectSquares()
    assert ds.count([0, 0]) == 0


def test_detect_squares_symmetric_both_sides():
    ds = DetectSquares()
    ds.add([0, 2])
    ds.add([2, 0])
    ds.add([2, 2])
    ds.add([-2, 0])
    ds.add([-2, 2])

    assert ds.count([0, 0]) == 2


def test_detect_squares_one_sided_only():
    ds = DetectSquares()
    ds.add([1, 4])
    ds.add([4, 1])
    ds.add([4, 4])

    assert ds.count([1, 1]) == 1


def test_detect_squares_duplicate_frequency_multiplication():
    ds = DetectSquares()
    ds.add([1, 4])
    ds.add([1, 4])
    ds.add([1, 4])
    ds.add([4, 1])
    ds.add([4, 1])
    ds.add([4, 4])

    assert ds.count([1, 1]) == 6


def test_detect_squares_no_matching_x_coordinate():
    ds = DetectSquares()
    ds.add([5, 5])
    ds.add([5, 9])
    ds.add([9, 5])
    ds.add([9, 9])

    assert ds.count([100, 100]) == 0


def test_detect_squares_count_does_not_mutate_state():
    ds = DetectSquares()
    ds.add([0, 2])
    ds.add([2, 0])
    ds.add([2, 2])
    ds.add([-2, 0])
    ds.add([-2, 2])

    assert ds.count([0, 0]) == 2
    assert ds.count([0, 0]) == 2

    ds.add([0, 2])
    assert ds.count([0, 0]) == 4
