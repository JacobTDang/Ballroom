from solution import k_closest


def normalize(points):
    return sorted(points)


def test_k_closest_case_1():
    assert normalize(k_closest([[1, 3], [-2, 2]], 1)) == normalize([[-2, 2]])


def test_k_closest_case_2():
    assert normalize(k_closest([[3, 3], [5, -1], [-2, 4]], 2)) == normalize([[3, 3], [-2, 4]])


def test_k_closest_case_3():
    assert normalize(k_closest([[0, 1], [1, 0]], 2)) == normalize([[0, 1], [1, 0]])


def test_k_closest_case_4():
    assert normalize(k_closest([[1, 1], [2, 2], [3, 3]], 3)) == normalize(
        [[1, 1], [2, 2], [3, 3]]
    )


def test_k_closest_case_5():
    assert normalize(k_closest([[-5, 4], [-6, -1], [3, 6], [2, -2]], 2)) == normalize(
        [[2, -2], [-6, -1]]
    )
