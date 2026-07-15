from solution import subsets_with_dup


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_subsets_with_dup_case_1():
    assert normalize(subsets_with_dup([1, 2, 2])) == normalize(
        [[], [1], [1, 2], [1, 2, 2], [2], [2, 2]]
    )


def test_subsets_with_dup_case_2():
    assert normalize(subsets_with_dup([0])) == normalize([[], [0]])


def test_subsets_with_dup_case_3():
    assert normalize(subsets_with_dup([4, 4, 4, 1, 4])) == normalize(
        [
            [],
            [1],
            [1, 4],
            [1, 4, 4],
            [1, 4, 4, 4],
            [1, 4, 4, 4, 4],
            [4],
            [4, 4],
            [4, 4, 4],
            [4, 4, 4, 4],
        ]
    )
