from solution import subsets


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_subsets_case_1():
    assert normalize(subsets([1, 2, 3])) == normalize(
        [[], [1], [2], [1, 2], [3], [1, 3], [2, 3], [1, 2, 3]]
    )


def test_subsets_case_2():
    assert normalize(subsets([0])) == normalize([[], [0]])


def test_subsets_case_3():
    assert normalize(subsets([1, 2])) == normalize([[], [1], [2], [1, 2]])


def test_subsets_case_4():
    assert normalize(subsets([-1, 1])) == normalize([[], [-1], [1], [-1, 1]])
