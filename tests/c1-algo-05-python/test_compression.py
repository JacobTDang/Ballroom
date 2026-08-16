from solution import compress


def test_example_basic_run():
    assert compress("aaab") == "a3b"


def test_example_mixed_runs():
    assert compress("aabcccccaaa") == "a2bc5a3"


def test_example_no_runs_keeps_original():
    assert compress("abc") == "abc"


def test_empty():
    assert compress("") == ""


def test_single_char():
    assert compress("a") == "a"


def test_all_equal():
    assert compress("aaaaa") == "a5"


def test_multi_digit_count():
    assert compress("a" * 12) == "a12"


def test_equal_length_encoding_keeps_original():
    assert compress("aabb") == "aabb"


def test_two_chars_not_strictly_shorter():
    assert compress("aa") == "aa"


def test_trailing_single_run_has_no_count():
    assert compress("baaa") == "ba3"
