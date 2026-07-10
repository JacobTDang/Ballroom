from solution import merge_triplets


def test_classic():
    triplets = [[2, 5, 3], [1, 8, 4], [1, 7, 5]]
    assert merge_triplets(triplets, [2, 7, 5]) is True


def test_classic_false():
    triplets = [[3, 4, 5], [4, 5, 6]]
    assert merge_triplets(triplets, [3, 2, 5]) is False


def test_single_exact():
    assert merge_triplets([[5, 5, 5]], [5, 5, 5]) is True


def test_all_poisoned():
    triplets = [[10, 1, 1], [1, 10, 1], [1, 1, 10]]
    assert merge_triplets(triplets, [5, 5, 5]) is False


def test_partial_match():
    triplets = [[2, 1, 1], [1, 2, 1]]
    assert merge_triplets(triplets, [2, 2, 2]) is False
