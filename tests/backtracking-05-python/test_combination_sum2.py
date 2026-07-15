from solution import combination_sum2


def normalize(lists):
    return sorted(sorted(l) for l in lists)


def test_combination_sum2_case_1():
    assert normalize(combination_sum2([10, 1, 2, 7, 6, 1, 5], 8)) == normalize(
        [[1, 1, 6], [1, 2, 5], [1, 7], [2, 6]]
    )


def test_combination_sum2_case_2():
    assert normalize(combination_sum2([2, 5, 2, 1, 2], 5)) == normalize([[1, 2, 2], [5]])


def test_combination_sum2_case_3():
    assert normalize(combination_sum2([1, 1, 1, 2, 2], 4)) == normalize([[1, 1, 2], [2, 2]])
