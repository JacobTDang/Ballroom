from solution import three_sum


def normalize(triplets):
    """Sort each triplet ascending, then sort the list of triplets —
    3Sum's valid outputs aren't uniquely ordered, so tests compare as
    sets rather than asserting an exact sequence."""
    return sorted(sorted(t) for t in triplets)


def test_three_sum_case_01():
    assert normalize(three_sum([-1, 0, 1, 2, -1, -4])) == normalize([[-1, -1, 2], [-1, 0, 1]])


def test_three_sum_case_02():
    assert normalize(three_sum([0, 1, 1])) == []


def test_three_sum_case_03():
    assert normalize(three_sum([0, 0, 0])) == [[0, 0, 0]]


def test_three_sum_case_04():
    assert normalize(three_sum([])) == []


def test_three_sum_case_05():
    assert normalize(three_sum([0, 0, 0, 0])) == [[0, 0, 0]]


def test_three_sum_case_06():
    assert normalize(three_sum([-2, 0, 1, 1, 2])) == normalize([[-2, 0, 2], [-2, 1, 1]])


def test_three_sum_case_07():
    assert normalize(three_sum([-3, -2, -1])) == []


def test_three_sum_case_08():
    assert normalize(three_sum([1, 2, 3])) == []


def test_three_sum_case_09():
    assert normalize(three_sum([1, -1])) == []


def test_three_sum_case_10():
    assert normalize(three_sum([3, -2, 1, 0, -1, -3, 2, -2, 0])) == normalize(
        [[-3, 0, 3], [-3, 1, 2], [-2, -1, 3], [-2, 0, 2], [-1, 0, 1]]
    )
