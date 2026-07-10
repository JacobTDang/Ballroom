from solution import subsets_with_dup


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_subsets_with_dup():
    assert normalize(subsets_with_dup([1, 2, 2])) == normalize(
        [[], [1], [1, 2], [1, 2, 2], [2], [2, 2]]
    )
    assert normalize(subsets_with_dup([0])) == normalize([[], [0]])
