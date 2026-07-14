from solution import three_sum


def normalize(triplets):
    """Sort each triplet ascending, then sort the list of triplets —
    3Sum's valid outputs aren't uniquely ordered, so tests compare as
    sets rather than asserting an exact sequence."""
    return sorted(sorted(t) for t in triplets)


def test_three_sum():
    assert normalize(three_sum([-1, 0, 1, 2, -1, -4])) == normalize([[-1, -1, 2], [-1, 0, 1]])
    assert normalize(three_sum([0, 1, 1])) == []
    assert normalize(three_sum([0, 0, 0])) == [[0, 0, 0]]
    assert normalize(three_sum([])) == []
    assert normalize(three_sum([0, 0, 0, 0])) == [[0, 0, 0]]
    assert normalize(three_sum([-2, 0, 1, 1, 2])) == normalize([[-2, 0, 2], [-2, 1, 1]])
    assert normalize(three_sum([-3, -2, -1])) == []
    assert normalize(three_sum([1, 2, 3])) == []
    assert normalize(three_sum([1, -1])) == []
    assert normalize(three_sum([3, -2, 1, 0, -1, -3, 2, -2, 0])) == normalize(
        [[-3, 0, 3], [-3, 1, 2], [-2, -1, 3], [-2, 0, 2], [-1, 0, 1]]
    )
