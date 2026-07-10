from solution import k_closest


def normalize(points):
    return sorted(points)


def test_k_closest():
    assert normalize(k_closest([[1, 3], [-2, 2]], 1)) == normalize([[-2, 2]])
    assert normalize(k_closest([[3, 3], [5, -1], [-2, 4]], 2)) == normalize([[3, 3], [-2, 4]])
    assert normalize(k_closest([[0, 1], [1, 0]], 2)) == normalize([[0, 1], [1, 0]])
