from solution import combination_sum


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_combination_sum():
    assert normalize(combination_sum([2, 3, 6, 7], 7)) == normalize([[2, 2, 3], [7]])
    assert normalize(combination_sum([2, 3, 5], 8)) == normalize(
        [[2, 2, 2, 2], [2, 3, 3], [3, 5]]
    )
    assert normalize(combination_sum([2], 1)) == []
