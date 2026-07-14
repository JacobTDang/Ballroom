from solution import permute


def normalize_exact(lists):
    return sorted(lists)


def test_permute():
    assert normalize_exact(permute([1, 2, 3])) == normalize_exact(
        [[1, 2, 3], [1, 3, 2], [2, 1, 3], [2, 3, 1], [3, 1, 2], [3, 2, 1]]
    )
    assert normalize_exact(permute([0, 1])) == normalize_exact([[0, 1], [1, 0]])
    assert normalize_exact(permute([1])) == [[1]]
    assert normalize_exact(permute([1, -1])) == normalize_exact([[1, -1], [-1, 1]])
