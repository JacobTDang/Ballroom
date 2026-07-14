from solution import num_distinct


def test_classic():
    assert num_distinct("rabbbit", "rabbit") == 3


def test_second_classic():
    assert num_distinct("babgbag", "bag") == 5


def test_exact_match():
    assert num_distinct("abc", "abc") == 1


def test_target_longer():
    assert num_distinct("abc", "abcd") == 0


def test_empty_target():
    assert num_distinct("abc", "") == 1


def test_empty_source():
    assert num_distinct("", "abc") == 0


def test_both_empty():
    assert num_distinct("", "") == 1


def test_repeated_chars_combinatoric():
    assert num_distinct("aaaa", "aa") == 6


def test_single_char_many_occurrences():
    assert num_distinct("aaa", "a") == 3
