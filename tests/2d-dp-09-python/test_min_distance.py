from solution import min_distance


def test_classic():
    assert min_distance("horse", "ros") == 3


def test_second_classic():
    assert min_distance("intention", "execution") == 5


def test_empty_first():
    assert min_distance("", "abc") == 3


def test_identical():
    assert min_distance("abc", "abc") == 0


def test_both_empty():
    assert min_distance("", "") == 0


def test_empty_second():
    assert min_distance("abc", "") == 3


def test_single_char_replace():
    assert min_distance("a", "b") == 1


def test_pure_insertion():
    assert min_distance("cat", "cats") == 1


def test_multi_op_mix():
    assert min_distance("sunday", "saturday") == 3
