from solution import pacific_atlantic


def normalize(lists):
    return sorted(lists)


def test_classic():
    heights = [
        [1, 2, 2, 3, 5],
        [3, 2, 3, 4, 4],
        [2, 4, 5, 3, 1],
        [6, 7, 1, 4, 5],
        [5, 1, 1, 2, 4],
    ]
    want = [[0, 4], [1, 3], [1, 4], [2, 2], [3, 0], [3, 1], [4, 0]]

    got = normalize(pacific_atlantic(heights))
    assert got == normalize(want)


def test_single_cell_flows_to_both():
    heights = [[1]]
    got = normalize(pacific_atlantic(heights))
    assert got == [[0, 0]]
