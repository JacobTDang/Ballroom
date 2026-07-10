from solution import subsets


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_subsets():
    assert normalize(subsets([1, 2, 3])) == normalize(
        [[], [1], [2], [1, 2], [3], [1, 3], [2, 3], [1, 2, 3]]
    )
    assert normalize(subsets([0])) == normalize([[], [0]])
