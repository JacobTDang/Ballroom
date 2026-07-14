from solution import missing_number


def test_classic():
    assert missing_number([3, 0, 1]) == 2


def test_missing_at_end():
    assert missing_number([0, 1]) == 2


def test_missing_at_start():
    assert missing_number([1, 2, 3]) == 0


def test_single_element():
    assert missing_number([0]) == 1


def test_missing_in_middle():
    assert missing_number([0, 1, 3, 4, 5]) == 2


def test_larger_shuffled():
    assert missing_number([9, 6, 4, 2, 3, 5, 7, 0, 1]) == 8
